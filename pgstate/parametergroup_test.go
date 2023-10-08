package pgstate_test

import (
	"encoding/gob"
	"testing"

	"github.com/jrottersman/lats/pgstate"
)

func Test_EncodeParameterGroups(t *testing.T) {
	pg := pgstate.ParameterGroup{}
	pgs := []pgstate.ParameterGroup{}
	pgs = append(pgs, pg)
	r := pgstate.EncodeParameterGroups(pgs)
	var result []pgstate.ParameterGroup
	dec := gob.NewDecoder(&r)
	err := dec.Decode(&result)
	if err != nil {
		t.Errorf("decode error: %s", err)
	}
	if len(pgs) != len(result) {
		t.Errorf("got %d expected %d", len(result), len(pgs))
	}
}

func Test_DecodeParameterGroups(t *testing.T) {
	pg := pgstate.ParameterGroup{}
	pgs := []pgstate.ParameterGroup{}
	pgs = append(pgs, pg)
	r := pgstate.EncodeParameterGroups(pgs)
	result := pgstate.DecodeParameterGroups(r)
	if len(pgs) != len(result) {
		t.Errorf("got %d expected %d", len(result), len(pgs))
	}
}
