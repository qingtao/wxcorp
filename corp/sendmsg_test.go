package corp

import "testing"

func TestNewSendMsgURL(t *testing.T) {
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
			name: "1",
			args: args{
				accessToken: "123456",
			},
			want: "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=123456",
		},
		{
			name: "2",
			args: args{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSendMsgURL(tt.args.url, tt.args.accessToken); got != tt.want {
				t.Errorf("NewSendMsgURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
