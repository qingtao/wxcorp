package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestJsAPITicketResponse_Validate(t *testing.T) {
	var s = `{"errcode":0,"errmsg":"ok","ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA","expires_in":7200}`
	var jsTicket JsAPITicketResponse
	json.Unmarshal([]byte(s), &jsTicket)
	tests := []struct {
		name    string
		ticket  *JsAPITicketResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			ticket:  nil,
			wantErr: true,
		},
		{
			name:    "2",
			ticket:  &jsTicket,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ticket.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("JsAPITicketResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetJsAPITicket(t *testing.T) {
	var s = `{"errcode":0,"errmsg":"ok","ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA","expires_in":7200}`
	var jsTicket JsAPITicketResponse
	json.Unmarshal([]byte(s), &jsTicket)

	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		switch accesstoken {
		case "wantOk":
		case "wantJSONErr":
			s = `{"errcode":0,"errmsg":VbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA","expires_in":7200}`
		default:
			s = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, s)
	}))
	type args struct {
		url         string
		accessToken string
		typ         string
	}
	tests := []struct {
		name       string
		args       args
		wantTicket *JsAPITicketResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:         "http://127.0.0.2:8080",
				accessToken: "a",
			},
			wantTicket: nil,
			wantErr:    true,
		},
		{
			name: "2",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
			},
			wantTicket: &jsTicket,
			wantErr:    false,
		},
		{
			name: "3",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
			},
			wantTicket: nil,
			wantErr:    true,
		},
		{
			name: "4",
			args: args{
				accessToken: "",
			},
			wantTicket: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTicket, err := GetJsAPITicket(tt.args.url, tt.args.accessToken, tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJsAPITicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTicket, tt.wantTicket) {
				t.Errorf("GetJsAPITicket() = %v, want %v", gotTicket, tt.wantTicket)
			}
		})
	}
}
