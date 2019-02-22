package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewGetUserURL(t *testing.T) {
	type args struct {
		url         string
		accessToken string
		userid      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{},
			want: "",
		},
		{
			name: "2",
			args: args{accessToken: "123456", userid: "111"},
			want: "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=123456&userid=111",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetUserURL(tt.args.url, tt.args.accessToken, tt.args.userid); got != tt.want {
				t.Errorf("NewGetUserURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGetSimpleListURL(t *testing.T) {
	type args struct {
		url          string
		typ          string
		accessToken  string
		departmentID int
		fetchChild   int
		status       []int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{},
		},
		{
			name: "11",
			args: args{accessToken: "123456", departmentID: 1},
			want: `https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=123456&department_id=1`,
		},
		{
			name: "2",
			args: args{typ: "simple", departmentID: 0},
		},
		{
			name: "3",
			args: args{typ: "simple", accessToken: "123456", departmentID: 1, fetchChild: 1, status: []int{1}},
			want: "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=123456&department_id=1&fetch_child=1&status=1",
		},
		{
			name: "4",
			args: args{typ: "detail", accessToken: "123456", departmentID: 1, fetchChild: 1},
			want: "https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=123456&department_id=1&fetch_child=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetUserListURL(tt.args.url, tt.args.typ, tt.args.accessToken, tt.args.departmentID, tt.args.fetchChild, tt.args.status); got != tt.want {
				t.Errorf("NewGetSimpleListURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDuplicateInt(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{a: []int{1, 1, 2, 3, 4, 8, 3, 2, 1}},
			want: []int{1, 2, 3, 4, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceIntRemoveDuplicate(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicateInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		res     *UserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			res:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.res.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UserResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserListResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		res     *UserListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			res:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.res.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UserListResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	var s = `{"errcode":0,"errmsg":"ok","userid":"zhangsan","name":"李四","department":[1,2],"order":[1,2],"position":"后台工程师","mobile":"15913215421","gender":"1","email":"zhangsan@gzdev.com","is_leader_in_dept":[1,0],"avatar":"http://wx.qlogo.cn/mmopen/ajNVdqHZLLA3WJ6DSZUfiakYe37PKnQhBIeOQBO4czqrnZDS79FH5Wm5m4X69TBicnHFlhiafvDwklOpZeXYQQ2icg/0","telephone":"020-123456","enable":1,"alias":"jackzhang","extattr":{"attrs":[{"type":0,"name":"文本名称","text":{"value":"文本"}},{"type":1,"name":"网页名称","web":{"url":"http://www.test.com","title":"标题"}}]},"status":1,"qr_code":"https://open.work.weixin.qq.com/wwopen/userQRCode?vcode=xxx","external_position":"产品经理","external_profile":{"external_corp_name":"企业简称","external_attr":[{"type":0,"name":"文本名称","text":{"value":"文本"}},{"type":1,"name":"网页名称","web":{"url":"http://www.test.com","title":"标题"}},{"type":2,"name":"测试app","miniprogram":{"appid":"wx8bd80126147df384","pagepath":"/index","title":"my miniprogram"}}]}}`
	var user UserResponse
	json.Unmarshal([]byte(s), &user)
	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		res := s
		switch accesstoken {
		case "wantOk":
		case "wantJSONErr":
			res = `"errcode":0,"errmsg":"ok","userid":"zhangsan"`
		default:
			res = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, res)
	}))
	type args struct {
		url         string
		accessToken string
		userid      string
	}
	tests := []struct {
		name     string
		args     args
		wantUser *UserResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
				userid:      "abc",
			},
			wantUser: &user,
			wantErr:  false,
		},
		{
			name: "2",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
				userid:      "abc",
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "3",
			args: args{
				url:         "",
				accessToken: "",
				userid:      "abc",
			},
			wantUser: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := GetUser(tt.args.url, tt.args.accessToken, tt.args.userid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetUser() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestGetUserList(t *testing.T) {
	var s = `{"errcode":0,"errmsg":"ok","userlist":[{"userid":"zhangsan","name":"李四","department":[1,2],"order":[1,2],"position":"后台工程师","mobile":"15913215421","gender":"1","email":"zhangsan@gzdev.com","is_leader_in_dept":[1,0],"avatar":"http://wx.qlogo.cn/mmopen/ajNVdqHZLLA3WJ6DSZUfiakYe37PKnQhBIeOQBO4czqrnZDS79FH5Wm5m4X69TBicnHFlhiafvDwklOpZeXYQQ2icg/0","telephone":"020-123456","enable":1,"alias":"jackzhang","status":1,"extattr":{"attrs":[{"type":0,"name":"文本名称","text":{"value":"文本"}},{"type":1,"name":"网页名称","web":{"url":"http://www.test.com","title":"标题"}}]},"qr_code":"https://open.work.weixin.qq.com/wwopen/userQRCode?vcode=xxx","external_position":"产品经理","external_profile":{"external_corp_name":"企业简称","external_attr":[{"type":0,"name":"文本名称","text":{"value":"文本"}},{"type":1,"name":"网页名称","web":{"url":"http://www.test.com","title":"标题"}},{"type":2,"name":"测试app","miniprogram":{"appid":"wx8bd80126147df384","pagepath":"/index","title":"miniprogram"}}]}}]}`
	var res UserListResponse
	json.Unmarshal([]byte(s), &res)
	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		res := s
		switch accesstoken {
		case "wantOk":
		case "wantJSONErr":
			res = `"errcode":0,"errmsg":"ok","userid":"zhangsan"`
		default:
			res = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, res)
	}))
	type args struct {
		url         string
		typ         string
		accessToken string
		departmenID int
		fetchChild  int
		status      []int
	}
	tests := []struct {
		name    string
		args    args
		wantRes *UserListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:         ht.URL,
				typ:         "",
				accessToken: "wantOk",
				departmenID: 1,
				fetchChild:  0,
			},
			wantRes: &res,
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				url:         ht.URL,
				typ:         "simple",
				accessToken: "wantJSONErr",
				departmenID: 1,
				fetchChild:  0,
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				url:         ht.URL,
				typ:         "",
				accessToken: "",
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				url:         "http://127.0.0.1:8080",
				typ:         "",
				accessToken: "a",
				departmenID: 1,
				fetchChild:  0,
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := GetUserList(tt.args.url, tt.args.typ, tt.args.accessToken, tt.args.departmenID, tt.args.fetchChild, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetUserList() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
