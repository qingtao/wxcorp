package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
				accessToken: "a",
				userid:      "abc",
			},
			wantUser: &UserResponse{
				ErrCode: 40014,
				ErrMsg:  "invalid access_token",
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				userid: "abc",
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "5",
			args: args{
				url:         "http://127.0.0.1:8081",
				accessToken: "a",
				userid:      "a",
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

func TestExtText_Validate(t *testing.T) {
	type fields struct {
		Value string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			fields:  fields{Value: ""},
			wantErr: false,
		},
		{
			name:    "2",
			fields:  fields{Value: "1234567890123"},
			wantErr: true,
		},
		{
			name:    "2",
			fields:  fields{Value: "Hello, 中国"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ExtText{
				Value: tt.fields.Value,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ExtText.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtWeb_Validate(t *testing.T) {
	type fields struct {
		Title string
		URL   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				Title: "",
				URL:   "",
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				Title: "test",
				URL:   "",
			},
			wantErr: true,
		},
		{
			name: "3",
			fields: fields{
				Title: "",
				URL:   "http://www.a.com",
			},
			wantErr: true,
		},
		{
			name: "4",
			fields: fields{
				Title: "test",
				URL:   "http://www.a.com",
			},
			wantErr: false,
		},
		{
			name: "5",
			fields: fields{
				Title: "testTestTestTest", //4xtest
				URL:   "http://www.a.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ExtWeb{
				Title: tt.fields.Title,
				URL:   tt.fields.URL,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ExtWeb.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtMiniprogram_Validate(t *testing.T) {
	type fields struct {
		Title    string
		AppID    string
		PagePath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				Title:    "",
				AppID:    "",
				PagePath: "http://www.a.com",
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				Title:    "test",
				AppID:    "",
				PagePath: "http://www.a.com",
			},
			wantErr: true,
		},
		{
			name: "3",
			fields: fields{
				Title:    "",
				AppID:    "wx8bd80126147daa",
				PagePath: "http://www.a.com",
			},
			wantErr: true,
		},
		{
			name: "4",
			fields: fields{
				Title:    "test",
				AppID:    "wx8bd80126147daa",
				PagePath: "http://www.a.com",
			},
			wantErr: false,
		},
		{
			name: "5",
			fields: fields{
				Title:    "testTestTestTest", //4xtest
				AppID:    "wx8bd80126147daa",
				PagePath: "http://www.a.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ExtMiniprogram{
				Title:    tt.fields.Title,
				AppID:    tt.fields.AppID,
				PagePath: tt.fields.PagePath,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ExtMiniprogram.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtAttrs_Validate(t *testing.T) {
	type fields struct {
		Attrs []ExtAttr
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "0",
			fields: fields{
				Attrs: []ExtAttr{
					ExtAttr{
						Type: 0,
						Name: "test",
						Text: ExtText{
							Value: "hello,世界!",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "1",
			fields: fields{
				Attrs: []ExtAttr{
					ExtAttr{
						Type: 1,
						Name: "test",
						Web: ExtWeb{
							Title: "Test",
							URL:   "",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "2",
			fields: fields{
				Attrs: []ExtAttr{
					ExtAttr{
						Type: 2,
						Name: "test",
						Miniprogram: ExtMiniprogram{
							Title:    "test",
							AppID:    "",
							PagePath: "",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "3",
			fields: fields{
				Attrs: []ExtAttr{
					ExtAttr{
						Type: 3,
						Name: "test",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ExtAttrs{
				Attrs: tt.fields.Attrs,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ExtAttrs.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_Validate(t *testing.T) {
	type fields struct {
		UserID           string
		NewUserID        string
		Name             string
		Department       []int
		Order            []int
		IsLeaderInDept   []int
		Position         string
		Mobile           string
		Gender           string
		Enable           int
		Email            string
		Avatar           string
		AvatarMediaID    string
		Telephone        string
		Alias            string
		Status           int
		ExtAttr          *ExtAttrs
		ToInvite         bool
		ExternalPosition string
		ExternalProfile  *ExternalProfile
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				Enable: -1,
			},
			wantErr: true,
		},
		{
			name: "2",
			fields: fields{
				Enable: 1,
			},
			wantErr: true,
		},
		{
			name: "3",
			fields: fields{
				Enable: 1,
				Name:   strings.Repeat("成员名称长度为1~64个utf8字符", 4),
			},
			wantErr: true,
		},
		{
			name: "4",
			fields: fields{
				Enable:     1,
				Name:       "test",
				Department: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			},
			wantErr: true,
		},
		{
			name: "5",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1},
			},
			wantErr: true,
		},
		{
			name: "6",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1},
			},
			wantErr: true,
		},
		{
			name: "7",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{1000},
			},
			wantErr: true,
		},
		{
			name: "7-1",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 2},
				Order:          []int{1000},
			},
			wantErr: true,
		},
		{
			name: "8",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{1 << 32, 0, 0},
			},
			wantErr: true,
		},
		{
			name: "9",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{0, 0, 0},
				Gender:         "3",
			},
			wantErr: true,
		},
		{
			name: "10",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{0, 0, 0},
				Gender:         "1",
				Position:       strings.Repeat("12345678", 17),
			},
			wantErr: true,
		},
		{
			name: "11",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{0, 0, 0},
				Gender:         "1",
				Email:          "www",
			},
			wantErr: true,
		},
		{
			name: "11-1",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{0, 0, 0},
				Gender:         "1",
				Email:          "@www",
			},
			wantErr: true,
		},
		{
			name: "11-2",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{0, 0, 0},
				Gender:         "1",
				Email:          "www@",
			},
			wantErr: true,
		},
		{
			name: "11-3",
			fields: fields{
				Enable:         1,
				Name:           "test",
				Department:     []int{1, 2, 3},
				IsLeaderInDept: []int{1, 0, 0},
				Order:          []int{0, 0, 0},
				Gender:         "1",
				Email:          strings.Repeat("ab@c.com", 9),
			},
			wantErr: true,
		},
		{
			name: "12",
			fields: fields{
				Enable:           1,
				Name:             "test",
				Department:       []int{1, 2, 3},
				IsLeaderInDept:   []int{1, 0, 0},
				Order:            []int{0, 0, 0},
				Gender:           "1",
				Email:            "ab@c.com",
				ExternalPosition: "1234567890123",
			},
			wantErr: true,
		},
		{
			name: "13",
			fields: fields{
				Enable:           1,
				Name:             "test",
				Department:       []int{1, 2, 3},
				IsLeaderInDept:   []int{1, 0, 0},
				Order:            []int{0, 0, 0},
				Gender:           "1",
				Email:            "ab@c.com",
				ExternalPosition: "中国abc",
			},
			wantErr: true,
		},
		{
			name: "14",
			fields: fields{
				Enable:           1,
				Name:             "test",
				Department:       []int{1, 2, 3},
				IsLeaderInDept:   []int{1, 0, 0},
				Order:            []int{0, 0, 0},
				Gender:           "1",
				Email:            "ab@c.com",
				ExternalPosition: "资深数据专家",
				ExtAttr: &ExtAttrs{
					Attrs: []ExtAttr{
						ExtAttr{
							Type: 1,
							Name: "test",
							Web: ExtWeb{
								Title: "a",
								URL:   "",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "15",
			fields: fields{
				Enable:           1,
				Name:             "test",
				Department:       []int{1, 2, 3},
				IsLeaderInDept:   []int{1, 0, 0},
				Order:            []int{0, 0, 0},
				Gender:           "1",
				Email:            "ab@c.com",
				ExternalPosition: "资深数据专家",
				ExtAttr: &ExtAttrs{
					Attrs: []ExtAttr{
						ExtAttr{
							Type: 1,
							Name: "test",
							Web: ExtWeb{
								Title: "a",
								URL:   "http://www.example.com",
							},
						},
					},
				},
				ExternalProfile: &ExternalProfile{
					ExternalCoprName: "name",
					ExternalAttr: []ExtAttr{
						ExtAttr{
							Type: 1,
							Name: "test",
							Web: ExtWeb{
								Title: "",
								URL:   "http://www.example.com",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "16",
			fields: fields{
				Enable:           1,
				Name:             "test",
				Department:       []int{1, 2, 3},
				IsLeaderInDept:   []int{1, 0, 0},
				Order:            []int{0, 0, 0},
				Gender:           "1",
				Email:            "ab@c.com",
				ExternalPosition: "资深数据专家",
				ExtAttr: &ExtAttrs{
					Attrs: []ExtAttr{
						ExtAttr{
							Type: 1,
							Name: "test",
							Web: ExtWeb{
								Title: "a",
								URL:   "http://www.example.com",
							},
						},
					},
				},
				ExternalProfile: &ExternalProfile{
					ExternalCoprName: "name",
					ExternalAttr: []ExtAttr{
						ExtAttr{
							Type: 1,
							Name: "test",
							Web: ExtWeb{
								Title: "test",
								URL:   "http://www.example.com",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := User{
				UserID:           tt.fields.UserID,
				NewUserID:        tt.fields.NewUserID,
				Name:             tt.fields.Name,
				Department:       tt.fields.Department,
				Order:            tt.fields.Order,
				IsLeaderInDept:   tt.fields.IsLeaderInDept,
				Position:         tt.fields.Position,
				Mobile:           tt.fields.Mobile,
				Gender:           tt.fields.Gender,
				Enable:           tt.fields.Enable,
				Email:            tt.fields.Email,
				Avatar:           tt.fields.Avatar,
				AvatarMediaID:    tt.fields.AvatarMediaID,
				Telephone:        tt.fields.Telephone,
				Alias:            tt.fields.Alias,
				Status:           tt.fields.Status,
				ExtAttr:          tt.fields.ExtAttr,
				ToInvite:         tt.fields.ToInvite,
				ExternalPosition: tt.fields.ExternalPosition,
				ExternalProfile:  tt.fields.ExternalProfile,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExternalProfile_Validate(t *testing.T) {
	type fields struct {
		ExternalCoprName string
		ExternalAttr     []ExtAttr
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				ExternalCoprName: "name",
				ExternalAttr: []ExtAttr{
					ExtAttr{
						Type: 0,
						Name: "extName",
						Text: ExtText{Value: "hello"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ExternalProfile{
				ExternalCoprName: tt.fields.ExternalCoprName,
				ExternalAttr:     tt.fields.ExternalAttr,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ExternalProfile.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
