package state

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Object struct {
	Object  interface{}
	Order   int
	ObjType string
}

func NewObject(obj interface{}, order int, objtype string) Object {
	return Object{
		Object:  obj,
		Order:   order,
		ObjType: objtype,
	}
}

type Stack struct {
	Name                  string
	RestorationObjectName string
	Objects               map[int][]Object //int is the order in which we restore
}

func (s Stack) Encoder(filelocation string) (*bytes.Buffer, error) {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(s)
	if err != nil {
		log.Fatalf("Error encoding our snapshot: %s", err)
		return nil, err
	}
	return &encoder, nil
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
