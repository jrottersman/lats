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
	mu             sync.Mutex
	StateLocations []stateKV
}

func (s StateManager) updateState(name string, filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	kv := stateKV{
		Object:       name,
		FileLocation: filename,
	}
	s.StateLocations = append(s.StateLocations, kv)

	return nil
}

func InitState(f string) error {
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

func ReadState(filename string) (StateManager, error) {
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
