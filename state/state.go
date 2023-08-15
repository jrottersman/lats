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
	StateLocations []stateKV `json:"stateLocations"`
}

func (s *StateManager) UpdateState(name string, filename string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	kv := stateKV{
		Object:       name,
		FileLocation: filename,
	}
	s.StateLocations = append(s.StateLocations, kv)
}

func (s *StateManager) SyncState(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, err := json.Marshal(s.StateLocations)
	if err != nil {
		fmt.Printf("Error creating json: %s", err)
		return err
	}
	err = os.WriteFile(filename, m, 0644)
	if err != nil {
		fmt.Printf("Error writing file %s", err)
		return err
	}
	return nil
}

func (s *StateManager) GetFile(object string) interface{} {
	return nil
}

func InitState(filename string) error {
	initStr := []string{}
	m, err := json.Marshal(initStr)
	if err != nil {
		fmt.Printf("Error initing empty string to json %s", err)
		return err
	}
	err = os.WriteFile(filename, m, 0644)
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
	var mu sync.Mutex
	sm := StateManager{
		mu,
		s,
	}
	return sm, err
}
