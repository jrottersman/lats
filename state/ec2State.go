package state

import (
	"bytes"
	"encoding/gob"
	"log/slog"

	"github.com/jrottersman/lats/cmd"
)

func EncodeSecurityGroups(sg cmd.SecurityGroupOutput) bytes.Buffer {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)

	err := enc.Encode(sg)
	if err != nil {
		slog.Error("Error encoding our database", "error", err)
	}
	return encoder
}
