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

func TestNewCreateDepartmentURL(t *testing.T) {
	type args struct {
		url         string
		accessToken string
	}
	var (
		accessToken   = "a"
		wantCreateURL = `https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token=a`
		wantUpdateURL = `https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token=a`
		wantDeleteURL = `https://qyapi.weixin.qq.com/cgi-bin/department/delete?access_token=a`
	)

	t.Run("create", func(t *testing.T) {
		got := NewCreateDepartmentURL("", "")
		if got != "" {
			t.Errorf("NewCreateDepartmentURL() = %v, want %v", got, "")
		}
		got = NewCreateDepartmentURL("", accessToken)
		if got != wantCreateURL {
			t.Errorf("NewCreateDepartmentURL() = %v, want %v", got, wantCreateURL)
		}
	})
	t.Run("update", func(t *testing.T) {
		got := NewUpdateDepartmentURL("", "")
		if got != "" {
			t.Errorf("NewUpdateDepartmentURL() = %v, want %v", got, "")
		}
		got = NewUpdateDepartmentURL("", accessToken)
		if got != wantUpdateURL {
			t.Errorf("NewUpdateDepartmentURL() = %v, want %v", got, wantUpdateURL)
		}
	})
	t.Run("delete", func(t *testing.T) {
		got := NewDeleteDepartmentURL("", "")
		if got != "" {
			t.Errorf("NewDeleteDepartmentURL() = %v, want %v", got, "")
		}
		got = NewDeleteDepartmentURL("", accessToken)
		if got != wantDeleteURL {
			t.Errorf("NewDeleteDepartmentURL() = %v, want %v", got, wantDeleteURL)
		}
	})
}

func TestDeleteDepartment(t *testing.T) {
	var s = `{"errcode":0,"errmsg":"deleted"}`
	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		switch accesstoken {
		case "ok":
			fmt.Fprint(w, s)
		case "json_error":
			fmt.Fprint(w, `"errcode":0,"errmsg":"deleted"}`)
		default:
			fmt.Fprint(w, `{"errcode":1,"errmsg":"未知错误"}`)
		}
	}))
	type args struct {
		url         string
		accessToken string
		id          int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				url:         ht.URL,
				accessToken: "ok",
				id:          2,
			},
		},
		{
			name: "2",
			args: args{
				url:         ht.URL,
				accessToken: "json_error",
				id:          2,
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				url: ht.URL,
				id:  2,
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				url:         "http://127.0.0.1:8082",
				accessToken: "a",
				id:          2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteDepartment(tt.args.url, tt.args.accessToken, tt.args.id)
			if err != nil {
				if tt.wantErr {
					t.Log(err.Error())
				} else {
					t.Errorf("DeleteDepartment() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestChangeDepartment(t *testing.T) {
	var s = `{"name":"广州研发中心","parentid":1,"order":1,"id":2}`
	var dept Department
	json.Unmarshal([]byte(s), &dept)
	var res = `{"errcode":0,"errmsg":"created","id":2}`
	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ss string
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		switch accesstoken {
		case "ok":
			ss = res
		case "json_error":
			ss = `{"name":"广州研发中心",parentid":1,"order":1,"id":2}`
		default:
			ss = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, ss)
	}))
	type args struct {
		url         string
		accessToken string
		dept        *Department
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{url: ht.URL, accessToken: "ok", dept: &dept},
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{url: ht.URL, accessToken: "json_error", dept: &dept},
			wantErr: true,
		},
		{
			name:    "3",
			args:    args{url: "http://127.0.0.1:8088", accessToken: "a", dept: &dept},
			wantErr: true,
		},
		{
			name:    "4",
			args:    args{url: "", accessToken: "", dept: &dept},
			wantErr: true,
		},
		{
			name:    "5",
			args:    args{url: "", accessToken: "ok", dept: &Department{Name: "广州研发中心", ParentID: 0, Order: 1000000, ID: 2}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run("create_"+tt.name, func(t *testing.T) {
			if err := CreateDepartment(tt.args.url, tt.args.accessToken, tt.args.dept); (err != nil) != tt.wantErr {
				t.Errorf("CreateDepartment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		t.Run("update_"+tt.name, func(t *testing.T) {
			if err := UpdateDepartment(tt.args.url, tt.args.accessToken, tt.args.dept); (err != nil) != tt.wantErr {
				t.Errorf("CreateDepartment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDepartment_Validate(t *testing.T) {
	type fields struct {
		ID       int
		Name     string
		ParentID int
		Order    int
	}
	type args struct {
		action string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				ID:       1,
				Name:     "广州研发中心",
				ParentID: 0,
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				ID:       1,
				Name:     "广州研发中心-广州研发中心",
				ParentID: 0,
			},
			wantErr: true,
		},
		{
			name: "3",
			fields: fields{
				ID:       1,
				Name:     "广州研发中心",
				ParentID: 1,
			},
			wantErr: true,
		},
		{
			name: "4",
			fields: fields{
				ID:       1,
				Name:     "",
				ParentID: 0,
			},
			args:    args{action: "create"},
			wantErr: true,
		},
		{
			name: "5",
			fields: fields{
				ID:       1,
				Name:     "广州研发中心?:",
				ParentID: 0,
			},
			args:    args{action: "create"},
			wantErr: true,
		},
		{
			name: "6",
			fields: fields{
				ID:       1 << 32,
				Name:     "广州研发中心",
				ParentID: 0,
			},
			args:    args{action: "create"},
			wantErr: true,
		},
		{
			name: "7",
			fields: fields{
				ID:       32,
				Name:     "广州研发中心",
				ParentID: 1 << 32,
			},
			args:    args{action: "create"},
			wantErr: true,
		},
		{
			name: "8",
			fields: fields{
				ID:       32,
				Name:     "广州研发中心",
				ParentID: 1,
				Order:    1 << 32,
			},
			args:    args{action: "create"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Department{
				ID:       tt.fields.ID,
				Name:     tt.fields.Name,
				ParentID: tt.fields.ParentID,
				Order:    tt.fields.Order,
			}
			if err := a.Validate(tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("Department.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeDepartmentResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		res     *ChangeDepartmentResponse
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
				t.Errorf("ChangeDepartmentResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
