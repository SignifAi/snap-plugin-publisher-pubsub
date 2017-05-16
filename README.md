[![CircleCI](https://img.shields.io/circleci/project/github/RedSparr0w/node-csgo-parser.svg)](https://circleci.com/gh/SignifAi/snap-plugin-publisher-pubsub)
[![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)](https://github.com/SignifAi/snap-plugin-publisher-pubsub/blob/master/LICENSE)

# snap-plugin-publisher-pubsub
Snap-Telematry Plugin for Google Cloud Pub/Sub.

Originally designed to support SignifAI.io pusblisher requirements over Google PubSub but can also be used to publish other generic metrics from the framework to Google PubSub service.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Community Support](#community-support)
3. [Contributing](#contributing)
4. [License](#license-and-authors)
5. [Acknowledgements](#acknowledgements)

## Getting Started
This plug-in was created to push metric events into SignifAI.io using Google Cloud PubSub. However, we thought it will be useful to create a generic plug-in for Google PubSub as much as possible. You can use it without any requirement for a SignifAI account.
The main reason it's useful is to be able and push metric events under HIPAA regulations. 
### System Requirements 
* [golang 1.8+](https://golang.org/dl/) (needed only for building)
  Context is in stdlib from 1.7.

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
#### Download Snap binary:
For Snap, check [here](https://github.com/intelsdi-x/snap/releases).


#### To build the plugin binary:
Fork https://github.com/SignifAi/snap-plugin-publisher-pubsub

Clone repo into `$GOPATH/src/github.com/SignifAi/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-pubsub.git
```

build:
  ```make```

Note: You can also change your local grpc to version4 (found in
intelsdi-*plugin/v1/rpc/*.pb.go)

testing:
  For full integration testing you'll need google cloud SDK so we can
use the pubsub emulator locally.

  You can download && install here: https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-149.0.0-darwin-x86_64.tar.gz

  ```make test```

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

#### Load the Plugin
Once the framework is up and running, you can load the plugin.
```
$ snaptel plugin load snap-plugin-publisher-pubsub
Plugin loaded
Name: pubsub
Version: 1
Type: publisher
Signed: false
Loaded Time: Sat, 18 Mar 2017 13:28:45 PDT
```

#### Task File
You need to create or update a task file to use the signafai publisher
plugin. We have provided an example, __tasks/signifai.yaml_ shown below. In
our example, we utilize the psutil collector so we have some data to
work with. There are three (3) configuration settings you can use.

Setting|Description|Required?|
|-------|-----------|---------|
|service_key|Google Service Account Key.|Yes|
|serialization|json or msgpack|choose one of them.|
|project_id|Your Google Cloud account|Yes|

```
---
  version: 1
  schedule:
    type: "simple"
    interval: "5s"
  max-failures: 10
  workflow:
    collect:
      config:
      metrics:
        /intel/psutil/load/load1: {} 
        /intel/psutil/load/load15: {}
        /intel/psutil/load/load5: {}
        /intel/psutil/vm/available: {}
        /intel/psutil/vm/free: {}
        /intel/psutil/vm/used: {}
      publish:
        - plugin_name: "gpubsub-publisher"
          config:
            host: "my.host"
            project_id: "prod-key"
            event_source: "Snap"
            serialization: "msgpack"
            service_key: '{
  "type": "service_account",
  "project_id": "prod-key",
  "private_key_id": "private_key_id",
  "private_key": "-----BEGIN PRIVATE KEY-----\nsome key goes here...n-----END PRIVATE KEY-----\n",
  "client_email": "sigtest@prod-1033.iam.gserviceaccount.com",
  "client_id": "115369188776538168476",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://accounts.google.com/o/oauth2/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/sigtest%40prod-1033.iam.gserviceaccount.com"
}'
```

Once the task file has been created, you can create and watch the task.
```
$ snaptel task create -t tasks/signafai.yaml
Using task manifest to create task
Task created
ID: 72869b36-def6-47c4-9db2-822f93bb9d1f
Name: Task-72869b36-def6-47c4-9db2-822f93bb9d1f
State: Running

$ snaptel task list
ID                                       NAME
STATE     ...
72869b36-def6-47c4-9db2-822f93bb9d1f
Task-72869b36-def6-47c4-9db2-822f93bb9d1f    Running   ...
```

## Community Support
Open an issue for any questions or request. We will do our best to support it.

## Contributing - We love contributions!

There's more than one way to give back, from examples to blogs to code updates. We accept all suggestions using pull requests.

## License
Released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@SignifAi](https://github.com/SignifAi/)
* Info: www.signifai.io
