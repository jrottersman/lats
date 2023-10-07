package state

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"
)

const LoneInstance = "SingleRDSInstance"
const Cluster = "RDSCluster"
const Instance = "ClusterInstance"
const DBParameterGroup = "DBParameterGroup"
const DBClusterParameterGroup = "DBClusterParameterGroup"

type Object struct {
	FileName string
	Order    int
	ObjType  string
}

// ReadObject Read the file for the object
func (o Object) ReadObject() interface{} {
	dat, err := os.ReadFile(o.FileName)
	if err != nil {
		fmt.Printf("error reading object file %s", err)
	}
	buf := bytes.NewBuffer(dat)
	switch o.ObjType {
	case LoneInstance:
		return DecodeRestoreDBInstanceFromDBSnapshotInput(*buf)
	case Cluster:
		return DecodeRestoreDBClusterFromSnapshotInput(*buf)
	case Instance:
		return DecodeCreateDBInstanceInput(*buf)
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
		log.Fatalf("Error encoding our stack: %s", err)
		return nil, err
	}
	return &encoder, nil
}

func (s Stack) Write(filename string) error {
	b, err := s.Encoder()
	if err != nil {
		log.Fatalf("Error creating bytes %s", err)
		return err
	}
	_, err = WriteOutput(filename, *b)
	if err != nil {
		log.Fatalf("error wrting output %s", err)
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
		fmt.Printf("Error reading the file %s", err)
		return nil, err
	}
	buf := bytes.NewBuffer(f)
	var stack Stack
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&stack)
	if err != nil {
		fmt.Printf("error decoding the gob %s", err)
		return nil, err
	}
	return &stack, nil
}

func DeleteStack(filename string) error {
	err := os.Remove(filename)
	return err
}
