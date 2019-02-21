package corp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetEchoStr(t *testing.T) {
	type args struct {
		corpid         string
		token          string
		encodingAESKey string
		s              string
	}
	tests := []struct {
		name        string
		args        args
		wantEchostr []byte
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				corpid:         "ww6a49152bad10fa40",
				s:              "/contact?msg_signature=e94e4549349cfaa88058d425e990237f06c01c0e&timestamp=1548646136&nonce=1548381701&echostr=DrdwnEhttJWlbTV1eFRz49NetjDXCrwPd3q%2BenXvHD4YBq6wEB4CmC21thyAy2fxw55D5L38YKMF55chC5MCQA%3D%3D",
				token:          "qbw2JZV581j",
				encodingAESKey: "G5pMmqEIYdO8qneyK3fdxdbizm4f2noJ8t0MdhT97iF",
			},
			wantEchostr: []byte("6467373778033604605"),
			wantErr:     false,
		},
		{
			name: "2",
			args: args{
				corpid:         "wx5823bf96d3bd56c7",
				s:              "/cgi-bin/wxpush?msg_signature=5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3&timestamp=1409659589&nonce=263014780&echostr=P9nAzCzyDtyTWESHep1vC5X9xho%2FqYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp%2B4RPcs8TgAE7OaBO%2BFZXvnaqQ%3D%3D",
				token:          "qDG6eK",
				encodingAESKey: "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C",
			},
			wantEchostr: nil,
			wantErr:     true,
		},
		{
			name: "3",
			args: args{
				s: "://a.b.com",
			},
			wantEchostr: nil,
			wantErr:     true,
		},
		{
			name: "4",
			args: args{
				corpid:         "wx5823bf96d3bd56c7",
				s:              "/cgi-bin/wxpush?timestamp=1409659589&nonce=263014780&echostr=P9nAzCzyDtyTWESHep1vC5X9xho%2FqYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp%2B4RPcs8TgAE7OaBO%2BFZXvnaqQ%3D%3D",
				token:          "qDG6eK",
				encodingAESKey: "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C",
			},
			wantEchostr: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEchostr, err := GetEchoStr(tt.args.corpid, tt.args.token, tt.args.encodingAESKey, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEchoStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEchostr, tt.wantEchostr) {
				fmt.Printf("%s\n", gotEchostr)
				t.Errorf("GetEchoStr() = %v, want %v", gotEchostr, tt.wantEchostr)
			}
		})
	}
}

func TestAccessTokenResponse_Validate(t *testing.T) {
	tests := []struct {
		name     string
		response *AccessTokenResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			response: &AccessTokenResponse{
				ErrCode: -1,
				ErrMsg:  "系统繁忙 服务器暂不可用",
			},
			wantErr: true,
		},
		{
			name:    "2",
			wantErr: true,
		},
		{
			name: "3",
			response: &AccessTokenResponse{
				ErrCode:     0,
				ErrMsg:      "ok",
				AccessToken: "123456",
				ExpiresIn:   12345678,
			},
			wantErr: false,
		},
		{
			name: "4",
			response: &AccessTokenResponse{
				ErrCode: -2,
				ErrMsg:  "系统繁忙 服务器暂不可用",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.response.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AccessTokenResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAccessTokenURL(t *testing.T) {
	type args struct {
		url    string
		corpid string
		secret string
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
				corpid: "1",
				secret: "2",
			},
			want: "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=1&corpsecret=2",
		},
		{
			name: "2",
			args: args{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAccessTokenURL(tt.args.url, tt.args.corpid, tt.args.secret); got != tt.want {
				t.Errorf("NewAccessTokenURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var s string
		corpid := r.FormValue("corpid")
		switch corpid {
		case "wantOk":
			s = `{"access_token":"accesstoken000001","expires_in":7200}`
		case "wantErrJson":
			s = `"access_token":"accesstoken000001","expires_in":7200}`
		case "readErr":
			s = "aaaaaaaaaaaaaaa"
		default:
			s = `{"errcode":40013,"errmsg":"invalid corpid"}`
		}
		fmt.Fprint(w, s)
	}))
	type args struct {
		url    string
		corpid string
		secret string
	}
	tests := []struct {
		name    string
		args    args
		wantRes *AccessTokenResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:    ts.URL,
				corpid: "40013",
				secret: "2",
			},
			wantRes: &AccessTokenResponse{
				ErrCode: 40013,
				ErrMsg:  "invalid corpid",
			},
			wantErr: true,
		},
		{
			name: "11",
			args: args{
				url:    ts.URL,
				corpid: "readErr",
				secret: "2",
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				url:    ts.URL,
				corpid: "wantOk",
				secret: "123456",
			},
			wantRes: &AccessTokenResponse{
				AccessToken: "accesstoken000001",
				ExpiresIn:   7200,
			},
			wantErr: false,
		},
		{
			name: "22",
			args: args{
				url:    ts.URL,
				corpid: "wantErrJson",
				secret: "123456",
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				corpid: "1",
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				url:    "https://127.0.0.1/reject",
				corpid: "1",
				secret: "2",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := GetAccessToken(tt.args.url, tt.args.corpid, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetAccessToken() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
