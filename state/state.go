package state

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"sync"
)

// THis whole approach might need some serious refactoring I should be using a map I think s

const SnapshotType = "RDSSnapshot"
const RdsInstanceType = "RDSInstance"
const KMSKeyType = "KMSKey"
const RdsClusterType = "RDSCluster"
const ClusterSnapshotType = "RDSClusterSnapshot"

// StateKV manages our state file and object location
type StateKV struct {
	Object       string `json:"object"`
	FileLocation string `json:"fileLocation"`
	ObjectType   string `json:"objectType"`
}

type StateManager struct {
	Mu             *sync.Mutex
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
		slog.Error("Error creating json", "error", err)
		return err
	}
	err = os.WriteFile(filename, m, 0644)
	if err != nil {
		slog.Error("Error writing file", "error", err)
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
				slog.Error("error reading the file", "error", err)
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
		slog.Error("Error initing empty string to json", "error", err)
		return err
	}
	err = os.WriteFile(filename, m, 0644)
	if err != nil {
		slog.Error("Error writing file", "error", err)
		return err
	}
	return nil
}

func ReadState(filename string) (StateManager, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Error reading the file", "error", err)
	}
	var s []StateKV
	err = json.Unmarshal(f, &s)
	if err != nil {
		slog.Error("Error reading the file", "error", err)
	}
	var mu sync.Mutex
	sm := StateManager{
		&mu,
		s,
	}
	return sm, err
}
