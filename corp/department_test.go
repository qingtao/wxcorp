package corp

import "testing"

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
