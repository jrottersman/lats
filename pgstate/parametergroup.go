package pgstate

import (
	"bytes"
	"encoding/gob"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type ParameterGroup struct {
	ParameterGroup        types.DBParameterGroup
	ClusterParameterGroup types.DBClusterParameterGroup
	Params                []types.Parameter
}

func EncodeParameterGroups(pgs []ParameterGroup) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(&pgs)
	if err != nil {
		slog.Error("Error encoding parameters", "Error", err)
	}
	return encoder
}

//DecodeParameterGroups Decodes the parameter Group
func DecodeParameterGroups(b bytes.Buffer) []ParameterGroup {
	var pg []ParameterGroup
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&pg)
	if err != nil {
		slog.Error("Error decoding parameters", "Error", err)
	}
	return pg
}
