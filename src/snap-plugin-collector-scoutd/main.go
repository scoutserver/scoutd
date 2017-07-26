package main

import (
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"snap-plugin-collector-scoutd/scout"
)

const (
	pluginName    = "scout-collector"
	pluginVersion = 1
)

func main() {
	plugin.StartCollector(scout.NewScoutCollector(), pluginName, pluginVersion)
}
