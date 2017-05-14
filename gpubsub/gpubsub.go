package gpubsub

import (
	"context"
	"errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"encoding/json"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/ugorji/go/codec"
	"log"

	"cloud.google.com/go/pubsub"
	"strings"
)

var MissingHostServiceApplication = errors.New("Your Configuration is Missing a Host, Service, or Application Field")

var MissingAuth = errors.New("You Configuration is Missing a Google Account Service Key")

// Publisher is a publisher to Google PubSub/SignifAi System
type Publisher struct {
	eventSource   string                   // Event Source of event - defaults to Snap
	host          string                   // host that is being collected from
	service       string                   // service that is being collected from
	application   string                   // application that is being collected from
	initialized   bool                     // indicates that we've initialized the plugin
	projectID     string                   // google cloud project id
	serialization string                   // serialization lib to use, valid options are {json, msgpack}
	client        *pubsub.Client           // google cloud pubsub client
	topics        map[string]*pubsub.Topic // map of topic to topic pointer
	ctx           context.Context          // google cloud context
}

func New() *Publisher {
	return new(Publisher)
}

// GetConfigPolicy returns the configuration Policy needed for using
// this plugin
//
// we have quite a few optional parameters here
func (p *Publisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{""}, "host", false)
	policy.AddNewStringRule([]string{""}, "service", false)
	policy.AddNewStringRule([]string{""}, "application", false)
	policy.AddNewStringRule([]string{""}, "event_source", true)
	policy.AddNewStringRule([]string{""}, "serialization", true)
	policy.AddNewStringRule([]string{""}, "service_key", true)

	return *policy, nil
}

// create topics setups the initial connection && creates the topics
// found in task
func (p *Publisher) createTopics(topics []string, service_key string) error {
	var err error
	p.ctx = context.Background()

	p.topics = make(map[string]*pubsub.Topic)

	ctx := context.Background()
	jwtConfig, err := google.JWTConfigFromJSON([]byte(service_key), pubsub.ScopePubSub)
	if err != nil {
		log.Println(err)
	}
	ts := jwtConfig.TokenSource(ctx)

	p.client, err = pubsub.NewClient(p.ctx, p.projectID, option.WithTokenSource(ts))
	if err != nil {
		return err
	}

	for i := 0; i < len(topics); i++ {
		topic, err := p.client.CreateTopic(p.ctx, topics[i])
		if err != nil {

			switch v := err.(type) {
			case *googleapi.Error:
				if v.Code == 409 {
					log.Printf("already created topic %v\n", topics[i])
					topic = p.client.Topic(topics[i])
				} else {
					return err
				}
			default:
				if grpc.Code(err) == codes.AlreadyExists {
					log.Printf("already created topic %v\n", topics[i])
					topic = p.client.Topic(topics[i])
				} else {
					return err
				}
			}

		}
		p.topics[topics[i]] = topic
	}

	return nil
}

// prob. want to refactor me
// the default Get* functions from plugin do assertations along w/nil
// chks
func (p *Publisher) setConfig(cfg plugin.Config, topics []string) error {

	if p.initialized {
		return nil
	}

	// mandatory
	project_id, err := cfg.GetString("project_id")
	if err != nil {
		log.Println(err)
		return err
	}
	p.projectID = project_id

	// mandatory
	event_source, err := cfg.GetString("event_source")
	if err != nil {
		log.Println(err)
		return err
	}
	p.eventSource = event_source

	host, err := cfg.GetString("host")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.host = host
	}

	service, err := cfg.GetString("service")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.service = service
	}

	application, err := cfg.GetString("application")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.application = application
	}

	// mandatory
	serialization, err := cfg.GetString("serialization")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.serialization = serialization
	}

	// mandatory
	service_key, err := cfg.GetString("service_key")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			return MissingAuth
		}
	}

	if service_key == "" {
		return MissingAuth
	}

	if p.host == "" && p.application == "" && p.service == "" {
		return MissingHostServiceApplication
	}

	err = p.createTopics(topics, service_key)
	if err != nil {
		return err
	}

	p.initialized = true

	return nil
}

func (p Publisher) extractTopics(mts []plugin.Metric) []string {
	var topics []string
	for _, m := range mts {
		var statics []string
		for _, element := range m.Namespace {
			if !element.IsDynamic() {
				statics = append(statics, element.Value)
			}
		}

		tname := strings.Join(statics, ".")
		// no wildcards allowed in namespace
		// maybe dynamics should go here?
		tname = strings.Replace(tname, ".*", "", -1)

		topics = append(topics, tname)
	}

	return topics
}

// Publish publishes snap metrics to Google PubSub
func (p *Publisher) Publish(mts []plugin.Metric, cfg plugin.Config) error {

	if !p.initialized {
		topics := p.extractTopics(mts)
		err := p.setConfig(cfg, topics)
		if err != nil {
			return err
		}
	}

	for _, m := range mts {

		var statics []string
		var attributes = make(map[string]interface{})
		for _, element := range m.Namespace {
			if element.IsDynamic() {
				attributes[element.Name] = element.Description
			} else {
				statics = append(statics, element.Value)
			}
		}

		tname := strings.Join(statics, ".")
		// no wildcards allowed in namespace
		// maybe dynamics should go here?
		tname = strings.Replace(tname, ".*", "", -1)

		o := Metric{
			EventSource: p.eventSource,
			EventType:   "metrics",
			Name:        tname,
			Value:       m.Data,
			Timestamp:   m.Timestamp.Unix(),
			Attributes:  attributes,
		}

		if p.host != "" {
			o.Host = p.host
		}

		if p.service != "" {
			o.Service = p.service
		}

		if p.application != "" {
			o.Application = p.application
		}

		var data = []byte{}
		if p.serialization == "msgpack" {
			var mh codec.MsgpackHandle
			enc := codec.NewEncoderBytes(&data, &mh)
			err := enc.Encode(o)
			if err != nil {
				return err
			}
		} else if p.serialization == "json" {
			var err error
			data, err = json.Marshal(o)
			if err != nil {
				return err
			}
		} else {
			return errors.New("no serialization set in task - although this should have been caught sooner")
		}

		p.topics[tname].Publish(p.ctx, &pubsub.Message{Data: data})

	}

	return nil
}
