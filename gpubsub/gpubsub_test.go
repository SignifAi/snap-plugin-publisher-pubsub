/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2017 SignifAI Inc
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gpubsub

import (
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"testing"
)

func TestValidConfig(t *testing.T) {
	p := New()

	config := make(plugin.Config)
	config["host"] = "my.local.host"
	config["project_id"] = "prod-1212"
	config["event_source"] = "Snap"
	config["serialization"] = "json"
	config["service_key"] = serviceKeyStub

	err := p.setConfig(config, []string{"topic1", "topic2"})
	if err != nil {
		t.Fatal(err)
	}

	if p.host != config["host"] {
		t.Fatalf("bad config, %v, %v", p.host, config["host"])
	}

	if p.projectID != config["project_id"] {
		t.Fatalf("bad config, %v, %v", p.projectID, config["project_id"])
	}

	if p.eventSource != config["event_source"] {
		t.Fatalf("bad config, %v, %v", p.eventSource, config["event_source"])
	}

	if p.serialization != config["serialization"] {
		t.Fatalf("bad config, %v, %v", p.serialization, config["serialization"])
	}

	if !p.initialized {
		t.Fatal("bad config, %v", p.initialized, true)
	}
}

func TestBadConfig(t *testing.T) {
	p := New()

	config := make(plugin.Config)
	config["api"] = "metrics"
	config["project_id"] = "prod-1212"
	config["event_source"] = "Snap"
	config["serializatione"] = "json"
	config["service_key"] = serviceKeyStub

	err := p.setConfig(config, []string{"topic1", "topic2"})
	if err != MissingHostServiceApplication {
		t.Fatal("mandatory field not erroring")
	}

	config["service"] = "my-webapp"

	config["service_key"] = ""

	err = p.setConfig(config, []string{"topic1", "topic2"})
	if err == nil || err != MissingAuth {
		t.Fatalf("mandatory field not erroring %v", err)
	}

	config["serialization"] = "msgpack"

	err = p.setConfig(config, []string{"topic1", "topic2"})
	if err != MissingAuth {
		t.Fatal("mandatory field not erroring")
	}

}
