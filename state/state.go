package state

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// StateKV manages our state file and object location
type stateKV struct {
	Object       string `json:"object"`
	FileLocation string `json:"fileLocation"`
}

type StateManager struct {
	mu sync.Mutex
	s  []stateKV
}

func initState(f string) error {
	initStr := []string{}
	m, err := json.Marshal(initStr)
	if err != nil {
		fmt.Printf("Error initing empty string to json %s", err)
		return err
	}
	err = os.WriteFile(f, m, 0644)
	if err != nil {
		fmt.Printf("Error writing file %s", err)
		return err
	}
	return nil
}

func readState(filename string) (StateManager, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading the file %s", err)
	}
	var s []stateKV
	err = json.Unmarshal(f, &s)
	if err != nil {
		fmt.Printf("Error reading the file %s", err)
	}
	var m sync.Mutex
	sm := StateManager{
		m,
		s,
	}
	return sm, err
}
