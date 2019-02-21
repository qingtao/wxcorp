package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewOAuth2RedirectURL(t *testing.T) {
	type args struct {
		wxurl     string
		appid     string
		returnTo  string
		targetURI string
		state     string
	}

	argsOk := args{
		appid:     "1002",
		returnTo:  "http://a.b.com",
		targetURI: "http://b.b.com",
		state:     "123456",
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: argsOk,
			want: `https://open.weixin.qq.com/connect/oauth2/authorize?appid=1002&redirect_uri=http%3A%2F%2Fb.b.com%3Freturn_to%3Dhttp%253A%252F%252Fa.b.com&response_type=code&scope=snsapi_base&state=123456#wechat_redirect`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOAuth2RedirectURL(tt.args.wxurl, tt.args.appid, tt.args.returnTo, tt.args.targetURI, tt.args.state); got != tt.want {
				t.Errorf("NewOAuth2RedirectURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGetUserInfoURL(t *testing.T) {
	type args struct {
		wxurl       string
		accessToken string
		code        string
	}

	argsOk := args{
		"", "a", "123456",
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: argsOk,
			want: "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=a&code=123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetUserInfoURL(tt.args.wxurl, tt.args.accessToken, tt.args.code); got != tt.want {
				t.Errorf("NewGetUserInfoURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInfoResponse_Validate(t *testing.T) {
	var s = `{"errcode":0,"errmsg":"ok","UserId":"USERID","DeviceId":"DEVICEID"}`
	var user UserInfoResponse
	json.Unmarshal([]byte(s), &user)
	tests := []struct {
		name    string
		user    *UserInfoResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			user:    nil,
			wantErr: true,
		},
		{
			name:    "2",
			user:    &user,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.user.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UserInfoResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserInfoWithCode(t *testing.T) {
	var s = `{"errcode":0,"errmsg":"ok","UserId":"USERID","DeviceId":"DEVICEID"}`
	var user UserInfoResponse
	json.Unmarshal([]byte(s), &user)
	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		switch accesstoken {
		case "wantOk":
		case "wantJSONErr":
			s = `{"errcode":1,"errmsg":"未知错误"`
		default:
			s = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, s)
	}))
	type args struct {
		url         string
		accessToken string
		code        string
	}
	tests := []struct {
		name    string
		args    args
		wantRes *UserInfoResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
				code:        "123456",
			},
			wantRes: &user,
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
				code:        "123456",
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				url:         "http://127.0.0.2:8080",
				accessToken: "a",
				code:        "123456",
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:    "4",
			args:    args{accessToken: ""},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := GetUserInfoWithCode(tt.args.url, tt.args.accessToken, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfoWithCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetUserInfoWithCode() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
