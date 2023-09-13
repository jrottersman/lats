package cmd

import (
	"os"
	"sync"
	"testing"

	"github.com/jrottersman/lats/state"
)

func TestGetState(t *testing.T) {
	// Generate state
	initState := state.StateManager{
		Mu: &sync.Mutex{},
	}
	initState.SyncState(".confState.json")
	defer os.Remove(".confState.json")
	// Generate config
	initConf := newConfig("us-east-1", "us-west-2")
	writeConfig(initConf, ".latsConfig.json")
	defer os.Remove(".latsConfig.json")

	config, _ := GetState()
	if config.MainRegion != "us-east-1" {
		t.Errorf("got %s expected us-east-1", config.MainRegion)
	}
}
