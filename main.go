package main

import (
	"github.com/SignifAi/snap-plugin-publisher-pubsub/gpubsub"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	pluginName    = "pubsub-publisher"
	pluginVersion = 1
)

func main() {
	plugin.StartPublisher(gpubsub.New(), pluginName, pluginVersion)
}
