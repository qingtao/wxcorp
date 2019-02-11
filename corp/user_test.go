package corp

import (
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
