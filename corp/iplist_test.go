package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestIPListResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ipList  *IPListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			ipList:  nil,
			wantErr: true,
		},
		{
			name: "2",
			ipList: &IPListResponse{
				ErrCode: 1,
				ErrMsg:  "未知错误",
				IPList:  nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ipList.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("IPListResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCallBackIPList(t *testing.T) {
	var res = `{"errcode":0,"errmsg":"ok","ip_list":["101.226.103.*","101.226.62.*"]}`
	var list IPListResponse
	json.Unmarshal([]byte(res), &list)
	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accessToken := r.FormValue("access_token")
		var s string
		switch {
		case accessToken == "wantOk":
			s = res
		case accessToken == "wantJSONErr":
			s = `{"errcode":0,"errmsg":"ok","ip_list":["101.226.103.*","101.226.62.*"}`
		default:
			s = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, s)
	}))
	type args struct {
		url         string
		accessToken string
	}
	tests := []struct {
		name       string
		args       args
		wantIPList *IPListResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				accessToken: "",
			},
			wantIPList: nil,
			wantErr:    true,
		},
		{
			name: "2",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
			},
			wantIPList: &list,
			wantErr:    false,
		},
		{
			name: "3",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
			},
			wantIPList: nil,
			wantErr:    true,
		},
		{
			name:       "4",
			args:       args{url: "http://127.0.0.2:8080", accessToken: "a"},
			wantIPList: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIPList, err := GetCallBackIPList(tt.args.url, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCallBackIPList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIPList, tt.wantIPList) {
				t.Errorf("GetCallBackIPList() = %v, want %v", gotIPList, tt.wantIPList)
			}
		})
	}
}

func TestNewIPListURL(t *testing.T) {
	type args struct {
		url         string
		accessToken string
	}
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:         "",
				accessToken: "",
			},
			wantS: "",
		},
		{
			name: "2",
			args: args{
				url:         "",
				accessToken: "access_token",
			},
			wantS: "https://qyapi.weixin.qq.com/cgi-bin/getcallbackip?access_token=access_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := NewIPListURL(tt.args.url, tt.args.accessToken); gotS != tt.wantS {
				t.Errorf("NewIPListURL() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
