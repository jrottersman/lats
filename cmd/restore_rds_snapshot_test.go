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
	ips := aws.PassedIPs{
		Port:        8000,
		Permissions: "127.0.0.1/32",
	}
	passedIps := []aws.PassedIPs{}
	passedIps = append(passedIps, ips)
	tests := []struct {
		name string
		args args
		want []aws.PassedIPs
	}{
		{
			name: "one",
			args: args{
				rules: []string{"127.0.0.1/32:8000"},
			},
			want: passedIps,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sgRuleConvert(tt.args.rules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sgRuleConvert() = %v, want %v", got, tt.want)
			}
		})
	}
}
