package helpers

import (
	"bytes"
	"testing"
)

func TestWriteOutput(t *testing.T) {
	type args struct {
		filename string
		b        bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "test", args: args{filename: "/tmp/test.txt", b: bytes.Buffer{}}, want: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteOutput(tt.args.filename, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
