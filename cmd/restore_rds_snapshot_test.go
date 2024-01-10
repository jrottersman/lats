package cmd

import (
	"reflect"
	"testing"

	"github.com/jrottersman/lats/aws"
)

func Test_sgRuleConvert(t *testing.T) {
	type args struct {
		rules []string
	}
	tests := []struct {
		name string
		args args
		want []aws.PassedIPs
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sgRuleConvert(tt.args.rules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sgRuleConvert() = %v, want %v", got, tt.want)
			}
		})
	}
}
