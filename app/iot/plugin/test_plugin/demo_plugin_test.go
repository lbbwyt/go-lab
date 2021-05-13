package main

import (
	"github.com/dullgiulio/pingo"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestTestPlugin_DeCode(t *testing.T) {
	// Make a new plugin from the executable we created. Connect to it via TCP
	p := pingo.NewPlugin("tcp", "/Users/mac/libaobao/github/private/go-lab/app/iot/plugin/test_plugin/test_plugin")
	// Actually start the plugin
	p.Start()
	// Remember to stop the plugin when done using it
	defer p.Stop()

	var (
		resp string
	)

	//rootCtx := context.Background()
	//ctx := context.WithValue(rootCtx, "request_Id", "123456")

	// Call a function from the object we created previously
	if err := p.Call("TestPlugin.DeCode", "123456", &resp); err != nil {
		log.Print(err)
	} else {
		log.Print(resp)
	}
}
