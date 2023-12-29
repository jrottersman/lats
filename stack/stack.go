package stack

import (
	"bytes"
	"encoding/gob"
	"log/slog"
	"os"
	"sort"

	"github.com/jrottersman/lats/helpers"
	"github.com/jrottersman/lats/pgstate"
	"github.com/jrottersman/lats/state"
)

const LoneInstance = "SingleRDSInstance"
const Cluster = "RDSCluster"
const Instance = "RDSInstance"
const DBParameterGroup = "DBParameterGroup"
const DBClusterParameterGroup = "DBClusterParameterGroup"
const OptionGroup = "OptionGroup"
const SecurityGroup = "SecurityGroup"

type Object struct {
	FileName string
	Order    int
	ObjType  string
}

// ReadObject Read the file for the object
func (o Object) ReadObject() interface{} {
	dat, err := os.ReadFile(o.FileName)
	slog.Info("filename is", "Filename", o.FileName)
	if err != nil {
		slog.Warn("error reading object file", "error", err)
	}
	buf := bytes.NewBuffer(dat)
	switch o.ObjType {
	case LoneInstance:
		return state.DecodeRestoreDBInstanceFromDBSnapshotInput(*buf)
	case Cluster:
		return state.DecodeRestoreDBClusterFromSnapshotInput(*buf)
	case Instance:
		return state.DecodeCreateDBInstanceInput(*buf)
	case DBParameterGroup:
		pgstate.DecodeParameterGroups(*buf)
	case DBClusterParameterGroup:
		pgstate.DecodeParameterGroups(*buf)
	case SecurityGroup:
		state.DecodeSecurityGroups(*buf)
	}
	return nil
}

func NewObject(filename string, order int, objtype string) Object {
	return Object{
		FileName: filename,
		Order:    order,
		ObjType:  objtype,
	}
}

type Stack struct {
	Name                  string //Name is the name of the stack
	RestorationObjectName string
	Objects               map[int][]Object //int is the order in which we restore
}

func (s Stack) Encoder() (*bytes.Buffer, error) {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(s)
	if err != nil {
		slog.Error("Error encoding our stack", "error", err)
		return nil, err
	}
	return &encoder, nil
}

func (s Stack) Write(filename string) error {
	b, err := s.Encoder()
	if err != nil {
		slog.Error("Error creating bytes", "error", err)
		return err
	}
	_, err = helpers.WriteOutput(filename, *b)
	if err != nil {
		slog.Error("error writing output", "error", err)
		return err
	}
	return nil
}

func (s Stack) SortStack() *[]int {
	keys := []int{}
	for k := range s.Objects {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return &keys
}

func NewStack(name string, restorationObjectName string, objects []Object) Stack {
	objs := make(map[int][]Object)
	for _, v := range objects {
		order := v.Order
		_, ok := objs[order]
		if !ok {
			objs[order] = []Object{v}
		} else {
			objs[order] = append(objs[order], v)
		}
	}

	return Stack{
		Name:                  name,
		RestorationObjectName: restorationObjectName,
		Objects:               objs,
	}
}

func ReadStack(filename string) (*Stack, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Error reading the stack", "error", err)
		return nil, err
	}
	buf := bytes.NewBuffer(f)
	var stack Stack
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&stack)
	if err != nil {
		slog.Error("Error Decoding Stack", "error", err)
		return nil, err
	}
	return &stack, nil
}

func DeleteStack(filename string) error {
	err := os.Remove(filename)
	return err
}
