package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTagListResponse_Validate(t *testing.T) {
	var s = `{
   "errcode": 0,
   "errmsg": "ok",
   "taglist":[
      {"tagid":1,"tagname":"a"},
      {"tagid":2,"tagname":"b"}
   ]
}`
	var list TagListResponse
	json.Unmarshal([]byte(s), &list)
	tests := []struct {
		name    string
		list    *TagListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			list:    &list,
			wantErr: false,
		},
		{
			name:    "2",
			list:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.list.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("TagListResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemberOfTagReponse_Validate(t *testing.T) {
	var s = `{
   "errcode": 0,
   "errmsg": "ok",
   "tagname": "乒乓球协会",
   "userlist": [
         {
             "userid": "zhangsan",
             "name": "李四"
         }
     ],
   "partylist": [2]
}`
	var member MemberOfTagReponse
	json.Unmarshal([]byte(s), &member)
	tests := []struct {
		name    string
		member  *MemberOfTagReponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			member:  &member,
			wantErr: false,
		},
		{
			name:    "2",
			member:  nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.member.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("MemberOfTagReponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewGetTagListURL(t *testing.T) {
	type args struct {
		url         string
		accessToken string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "2",
			args: args{},
			want: "",
		},
		{
			name: "2",
			args: args{
				accessToken: "a",
			},
			want: "https://qyapi.weixin.qq.com/cgi-bin/tag/list?access_token=a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetTagListURL(tt.args.url, tt.args.accessToken); got != tt.want {
				t.Errorf("NewGetTagListURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGetUserOfTagURL(t *testing.T) {
	type args struct {
		url         string
		accessToken string
		tagid       int
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
				tagid: 1,
			},
			want: "",
		},
		{
			name: "1",
			args: args{
				accessToken: "a",
				tagid:       1,
			},
			want: "https://qyapi.weixin.qq.com/cgi-bin/tag/get?access_token=a&tagid=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetUserOfTagURL(tt.args.url, tt.args.accessToken, tt.args.tagid); got != tt.want {
				t.Errorf("NewGetUserOfTagURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTagList(t *testing.T) {
	var s = `{
   "errcode": 0,
   "errmsg": "ok",
   "taglist":[
      {"tagid":1,"tagname":"a"},
      {"tagid":2,"tagname":"b"}
   ]
}`
	var tag TagListResponse
	json.Unmarshal([]byte(s), &tag)

	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		res := s
		switch accesstoken {
		case "wantOk":
		case "wantJSONErr":
			res = `{
   "errcode": 0,
   "errmsg": "ok",
   "taglist":[
      {"tagid":1,"tagname":"a"},
      "tagid":2,"tagname":"b"}
   ]
}`
		default:
			res = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, res)
	}))
	type args struct {
		url         string
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		wantRes *TagListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "wantOk",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
			},
			wantRes: &tag,
			wantErr: false,
		},
		{
			name: "wantJSONErr",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:    "1",
			args:    args{accessToken: ""},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "networkErr",
			args: args{
				url:         "http://127.0.0.1:8080",
				accessToken: "a",
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := GetTagList(tt.args.url, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetTagList() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGetMemberOfTag(t *testing.T) {
	var s = `{
   "errcode": 0,
   "errmsg": "ok",
   "tagname": "乒乓球协会",
   "userlist": [
         {
             "userid": "zhangsan",
             "name": "李四"
         }
     ],
   "partylist": [2]
}`
	var member MemberOfTagReponse
	json.Unmarshal([]byte(s), &member)

	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		res := s
		switch accesstoken {
		case "wantOk":
		case "wantJSONErr":
			res = `"errmsg": "ok",
   "tagname": "乒乓球协会",
   "userlist": [
         {
             "userid": "zhangsan",
             "name": "李四"
         }
     ],
   "partylist": [2]
}`
		default:
			res = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, res)
	}))
	type args struct {
		url         string
		accessToken string
		tagid       int
	}
	tests := []struct {
		name       string
		args       args
		wantMember *MemberOfTagReponse
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
			},
			wantMember: &member,
			wantErr:    false,
		},
		{
			name: "2",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
			},
			wantMember: nil,
			wantErr:    true,
		},
		{
			name: "3",
			args: args{
				url:         "",
				accessToken: "",
			},
			wantMember: nil,
			wantErr:    true,
		},
		{
			name: "4",
			args: args{
				url:         "http://127.0.0.1:8080",
				accessToken: "a",
			},
			wantMember: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMember, err := GetMemberOfTag(tt.args.url, tt.args.accessToken, tt.args.tagid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMemberOfTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMember, tt.wantMember) {
				t.Errorf("GetMemberOfTag() = %v, want %v", gotMember, tt.wantMember)
			}
		})
	}
}
