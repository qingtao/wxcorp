package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewGetDepartmentListURL(t *testing.T) {
	type args struct {
		url         string
		accessToken string
		id          int
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
				id:          1111,
			},
			want: "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=123456&id=1111",
		},
		{
			name: "2",
			args: args{
				accessToken: "",
				id:          1,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetDepartmentListURL(tt.args.url, tt.args.accessToken, tt.args.id); got != tt.want {
				t.Errorf("NewGetDepartmentListURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepartmentResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dept    *DepartmentResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			dept: &DepartmentResponse{
				ErrCode:    1,
				ErrMsg:     "未知错误",
				Department: []Department{},
			},
			wantErr: true,
		},
		{
			name:    "2",
			dept:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dept.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DepartmentResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDepartment(t *testing.T) {
	var res = `{"errcode":0,"errmsg":"ok","department":[{"id":2,"name":"广州研发中心","parentid":1,"order":10},{"id":3,"name":"邮箱产品部","parentid":2,"order":40}]}`
	var deptRes DepartmentResponse
	json.Unmarshal([]byte(res), &deptRes)
	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		token := r.FormValue("access_token")
		// departid := r.FormValue("departmentid")

		var s string
		switch {
		case token == "wantOk":
			s = res
		case token == "wantJSONErr":
			s = `{"errcode":0"errmsg":"ok","department":[{"id":2,"name":"广州研发中心","parentid":1,"order":10},{"id":3,"name":"邮箱产品部","parentid":2,"order":40}]}`
		default:
			s = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprintf(w, "%s", s)
	}))
	type args struct {
		url         string
		accessToken string
		id          int
	}
	tests := []struct {
		name     string
		args     args
		wantDept *DepartmentResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "wantOk",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
				id:          1,
			},
			wantDept: &deptRes,
			wantErr:  false,
		},
		{
			name: "accessTokenEmpty",
			args: args{
				url:         ht.URL,
				accessToken: "",
			},
			wantDept: nil,
			wantErr:  true,
		},
		{
			name: "networkErr",
			args: args{
				url:         "http://127.0.0.1:8080",
				accessToken: "wantOk",
			},
			wantDept: nil,
			wantErr:  true,
		},
		{
			name: "wantJSONErr",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
			},
			wantDept: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDept, err := GetDepartment(tt.args.url, tt.args.accessToken, tt.args.id)
			t.Logf("%v", err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDepartment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDept, tt.wantDept) {
				t.Errorf("GetDepartment() = %v, want %v", gotDept, tt.wantDept)
			}
		})
	}
}
