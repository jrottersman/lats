package state

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type Object struct {
	FileName string
	Order    int
	ObjType  string
}

func NewObject(filename string, order int, objtype string) Object {
	return Object{
		FileName: filename,
		Order:    order,
		ObjType:  objtype,
	}
}

type Stack struct {
	Name                  string
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
