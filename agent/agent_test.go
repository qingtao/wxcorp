package agent

import (
	"reflect"
	"testing"
)

func TestNewAgent(t *testing.T) {
	type args struct {
		url            string
		corpid         string
		agentid        string
		secret         string
		encodingAESKey string
		token          string
	}
	tests := []struct {
		name string
		args args
		want *Agent
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				corpid:         "1",
				agentid:        "2",
				secret:         "3",
				encodingAESKey: "11",
				token:          "22",
			},
			want: &Agent{
				CorpID:         "1",
				AgentID:        "2",
				Secret:         "3",
				EncodingAESKey: "11",
				Token:          "22",
				ipList:         make(map[string]struct{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAgent(tt.args.corpid, tt.args.agentid, tt.args.secret, tt.args.encodingAESKey, tt.args.token)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
