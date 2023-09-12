package state

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/jrottersman/lats/helpers"
)

// THis whole approach might need some serious refactoring I should be using a map I think s

const SnapshotType = "RDSSnapshot"
const RdsInstanceType = "RDSInstance"
const KMSKeyType = "KMSKey"
const RdsClusterType = "RDSCluster"
const ClusterSnapshotType = "RDSClusterSnapshot"

type StackStore struct {
	Store []StackLookup `json:"store"`
}

type StackLookup struct {
	Name string `json:"name"`
	File string `json:"file"`
}

//CreateStackLookUp writes a stack to disk and creates a LookUp for that stack
func CreateStackLookUp(stack Stack, filename ...string) StackLookup {
	if filename == nil {
		filename[0] = *helpers.RandomStateFileName()
	}
	err := stack.Write(filename[0])
	if err != nil {
		fmt.Printf("error creating stack %s", err)
	}
	return StackLookup{
		Name: stack.Name,
		File: filename[0],
	}
}

//StackFiles stores the locations of our stacks
type StackFiles struct {
	Stacks []StackLookup
}

//Append to our StackFiles
func (sf *StackFiles) AppendStackLookup(sl StackLookup) {
	sf.Stacks = append(sf.Stacks, sl)
}

//GetStack get's a single stack from our stack list
func (sf *StackFiles) GetStack(name string) (*Stack, error) {
	for _, v := range sf.Stacks {
		if v.Name == name {
			stack, err := ReadStack(v.File)
			if err != nil {
				return nil, err
			}
			return stack, err
		}
	}
	return nil, fmt.Errorf("Stack with name %s doesn't exist", name)
}

func (sf *StackFiles) RemoveStack(name string) {
	for i, v := range sf.Stacks {
		if v.Name == name {
			DeleteStack(v.File)
			sf.Stacks = append(sf.Stacks[:i], sf.Stacks[i+1:]...)
		}
	}
}

// StateKV manages our state file and object location
type StateKV struct {
	Object       string `json:"object"`
	FileLocation string `json:"fileLocation"`
	ObjectType   string `json:"objectType"`
}

type StateManager struct {
	Mu             sync.Mutex
	StateLocations []StateKV `json:"stateLocations"`
}

func (s *StateManager) UpdateState(name string, filename string, ot string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	kv := StateKV{
		Object:       name,
		FileLocation: filename,
		ObjectType:   ot,
	}
	s.StateLocations = append(s.StateLocations, kv)
}

func (s *StateManager) SyncState(filename string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
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

func (s *StateManager) GetStateObject(object string) interface{} {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for i := range s.StateLocations {
		if s.StateLocations[i].Object == object {
			dat, err := os.ReadFile(s.StateLocations[i].FileLocation)
			if err != nil {
				fmt.Printf("error reading the file %s", err)
			}
			buf := bytes.NewBuffer(dat)
			switch objType := s.StateLocations[i].ObjectType; objType {
			case SnapshotType:
				snap := DecodeRDSSnapshotOutput(*buf)
				return snap
			case RdsInstanceType:
				instance := DecodeRDSDatabaseOutput(*buf)
				return instance
			case RdsClusterType:
				cluster := DecodeRDSClusterOutput(*buf)
				return cluster
			case KMSKeyType:
				key := DecodeKmsOutput(*buf)
				return key
			case ClusterSnapshotType:
				csnap := DecodeRDSClusterSnapshotOutput(*buf)
				return csnap
			}
		}
	}
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
	var s []StateKV
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
