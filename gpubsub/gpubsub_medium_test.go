package gpubsub

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	plugin "github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"os"
	"testing"
	"time"
)

var serviceKeyStub = `{
  "type": "service_account",
  "project_id": "prod-1212",
  "private_key_id": "deadbeefcafe",
  "private_key": "-----BEGIN PRIVATE KEY-----\nsome-long-key\n-----END PRIVATE KEY-----\n",
  "client_email": "sigtest@prod-1033.iam.gserviceaccount.com",
  "client_id": "33333333",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://accounts.google.com/o/oauth2/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/sigtest%40prod-1033.iam.gserviceaccount.com"
}`

func validConfig() plugin.Config {
	config := make(plugin.Config)
	config["api"] = "metrics"
	config["host"] = "my.local.host"
	config["project_id"] = "prod-1212"
	config["event_source"] = "Snap"
	config["serialization"] = "json"
	config["service_key"] = serviceKeyStub

	return config
}

func TestPubSubPublisher(t *testing.T) {

	// talk to the local emulator
	os.Setenv("PUBSUB_EMULATOR_HOST", "127.0.0.1:8321")

	ctx := context.Background()

	p := New()

	metrics := []plugin.Metric{
		plugin.Metric{
			Namespace: plugin.NewNamespace("x", "y", "z"),
			Config:    map[string]interface{}{"pw": "123aB"},
			Data:      3,
			Tags:      map[string]string{"hello": "world"},
			Unit:      "int",
			Timestamp: time.Now(),
		},
		plugin.Metric{
			Namespace: plugin.NewNamespace("bar").AddDynamicElement("domain_name", "Domain Name"),
			Config:    map[string]interface{}{"pw": "123aB"},
			Data:      3,
			Tags:      map[string]string{"hello": "world"},
			Unit:      "int",
			Timestamp: time.Now(),
		},
	}
	err := p.Publish(metrics, validConfig())
	if err != nil {
		t.Fatal(err)
	}

	p.topics["x.y.z"].Stop()
	p.topics["bar"].Stop()

	// Creates a client.
	client, err := pubsub.NewClient(ctx, "prod-1212")
	if err != nil {
		t.Fatal(err)
	}

	xyzTop := client.Topic("x.y.z")
	xyzSub, err := client.CreateSubscription(ctx, "xyzsub", xyzTop, 0, nil)
	if err != nil {
		if grpc.Code(err) == codes.AlreadyExists {
			log.Println("prob. prior run should not exist")
			xyzSub = client.Subscription("xyzsub")
		} else {
			t.Fatal(err)
		}
	}

	go func() {
		err = xyzSub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
			x := Metric{}
			json.Unmarshal(m.Data, &x)
			if x.EventSource != "Snap" {
				t.Fatal("event source doesn't match")
			}

			if x.EventType != "metrics" {
				t.Fatal("event type is not metrics")
			}

			if x.Host != "my.local.host" {
				t.Fatal("host field is not set correctly")
			}

			if x.Name != "x.y.z" {
				t.Fatal("this metric name is incorrect")
			}

			// https://github.com/intelsdi-x/snap-plugin-lib-go/blob/master/v1/plugin/metric.go#L71
			val, ok := x.Value.(float64)
			if !ok {
				t.Fatal("wrong type for value")
			}

			if val != 3 {
				t.Fatal("this metric value is incorrect")
			}

			m.Ack()
		})
		if err != nil {
			t.Fatal(err)
		}
	}()

	barTop := client.Topic("bar")
	barSub, err := client.CreateSubscription(ctx, "barsub", barTop, 0, nil)
	if err != nil {
		if grpc.Code(err) == codes.AlreadyExists {
			log.Println("prob. prior run should not exist")
			xyzSub = client.Subscription("barsub")
		} else {
			t.Fatal(err)
		}
	}

	go func() {
		err = barSub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
			res := Metric{}
			json.Unmarshal(m.Data, &res)
			if res.EventSource != "Snap" {
				t.Fatal("event source doesn't match")
			}

			if res.EventType != "metrics" {
				t.Fatal("event type is not metrics")
			}

			if res.Host != "my.local.host" {
				t.Fatal("host field is not set correctly")
			}

			if res.Name != "bar" {
				t.Fatal("this metric name is incorrect")
			}

			val, ok := res.Value.(float64)
			if !ok {
				t.Fatal("wrong type for value")
			}

			if val != 3 {
				t.Fatal("this metric value is incorrect")
			}

			m.Ack()
		})
		if err != nil {
			t.Fatal(err)
		}

	}()

	//	p.topics["x.y.z"].Stop()
	//	p.topics["bar"].Stop()

}
