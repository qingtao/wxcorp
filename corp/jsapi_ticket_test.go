package corp

import (
	"reflect"
	"testing"
)

func TestNewJsAPITicketURL(t *testing.T) {
	type args struct {
		url         string
		accessToken string
		typ         string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				accessToken: "123456",
			},
			want: "https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=123456",
		},
		{
			name: "2",
			args: args{
				accessToken: "123456",
				typ:         "agent_config",
			},
			want: "https://qyapi.weixin.qq.com/cgi-bin/ticket/get?access_token=123456&type=agent_config",
		},
		{
			name: "3",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsAPITicketURL(tt.args.url, tt.args.accessToken, tt.args.typ); got != tt.want {
				t.Errorf("NewJsAPITicketURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenSignature(t *testing.T) {
	type args struct {
		ticket    string
		noncestr  string
		url       string
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				ticket:    "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg",
				noncestr:  "Wm3WZYTPz0wzccnW",
				url:       "http://mp.weixin.qq.com?params=value",
				timestamp: 1414587457,
			},
			want: "0f9de62fce790f9a083d5c99e95740ceb90c27ed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genSignature(tt.args.ticket, tt.args.noncestr, tt.args.url, tt.args.timestamp); got != tt.want {
				t.Errorf("GenSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJsAPITicketSignature(t *testing.T) {
	type args struct {
		corpid    string
		agentid   string
		ticket    string
		noncestr  string
		url       string
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want *JsAPITicketSignature
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				corpid:    "1",
				agentid:   "2",
				ticket:    "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg",
				noncestr:  "Wm3WZYTPz0wzccnW",
				url:       "http://mp.weixin.qq.com?params=value",
				timestamp: 1414587457,
			},
			want: &JsAPITicketSignature{
				CorpID:    "1",
				AgentID:   "2",
				Timestamp: 1414587457,
				Signature: "0f9de62fce790f9a083d5c99e95740ceb90c27ed",
				NonceStr:  "Wm3WZYTPz0wzccnW",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsAPITicketSignature(tt.args.corpid, tt.args.agentid, tt.args.ticket, tt.args.noncestr, tt.args.url, tt.args.timestamp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsAPITicketSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
