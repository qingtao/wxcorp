package corp

import (
	"testing"
)

func TestNewTextReplyMsg(t *testing.T) {
	type args struct {
		toUserID   string
		fromCorpID string
		content    string
		createTime int64
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
				toUserID:   "toUser",
				fromCorpID: "fromUser",
				createTime: 1348831860,
				content:    "this is a test",
			},
			want: "<xml><ToUserName><![CDATA[toUser]]></ToUserName><FromUserName><![CDATA[fromUser]]></FromUserName><CreateTime>1348831860</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[this is a test]]></Content></xml>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextReplyMsg(tt.args.toUserID, tt.args.fromCorpID, tt.args.content, tt.args.createTime); got != tt.want {
				t.Errorf("NewTextReplyMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVoiceReplyMsg(t *testing.T) {
	type args struct {
		toUserID   string
		fromCorpID string
		mediaID    string
		createTime int64
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
				toUserID:   "toUser",
				fromCorpID: "fromUser",
				createTime: 1357290913,
				mediaID:    "media_id",
			},
			want: "<xml><ToUserName><![CDATA[toUser]]></ToUserName><FromUserName><![CDATA[fromUser]]></FromUserName><CreateTime>1357290913</CreateTime><MsgType><![CDATA[voice]]></MsgType><Voice><MediaId><![CDATA[media_id]]></MediaId></Voice></xml>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVoiceReplyMsg(tt.args.toUserID, tt.args.fromCorpID, tt.args.mediaID, tt.args.createTime); got != tt.want {
				t.Errorf("NewVoiceReplyMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewImageReplyMsg(t *testing.T) {
	type args struct {
		toUserID   string
		fromCorpID string
		mediaID    string
		createTime int64
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
				toUserID:   "toUser",
				fromCorpID: "fromUser",
				createTime: 1348831860,
				mediaID:    "media_id",
			},
			want: "<xml><ToUserName><![CDATA[toUser]]></ToUserName><FromUserName><![CDATA[fromUser]]></FromUserName><CreateTime>1348831860</CreateTime><MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[media_id]]></MediaId></Image></xml>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImageReplyMsg(tt.args.toUserID, tt.args.fromCorpID, tt.args.mediaID, tt.args.createTime); got != tt.want {
				t.Errorf("NewImageReplyMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVideoReplyMsg(t *testing.T) {
	type args struct {
		toUserID   string
		fromCorpID string
		mediaID    string
		title      string
		desc       string
		createTime int64
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
				toUserID:   "toUser",
				fromCorpID: "fromUser",
				createTime: 1357290913,
				mediaID:    "media_id",
				title:      "title",
				desc:       "description",
			},
			want: "<xml><ToUserName><![CDATA[toUser]]></ToUserName><FromUserName><![CDATA[fromUser]]></FromUserName><CreateTime>1357290913</CreateTime><MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[media_id]]></MediaId><Title><![CDATA[title]]></Title><Description><![CDATA[description]]></Description></Video></xml>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVideoReplyMsg(tt.args.toUserID, tt.args.fromCorpID, tt.args.mediaID, tt.args.title, tt.args.desc, tt.args.createTime); got != tt.want {
				t.Errorf("NewVideoReplyMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNewsItemMsg(t *testing.T) {
	type args struct {
		title  string
		desc   string
		picURL string
		url    string
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
				title:  "title1",
				desc:   "description1",
				picURL: "picurl",
				url:    "url",
			},
			want: "<item><Title><![CDATA[title1]]></Title><Description><![CDATA[description1]]></Description><PicUrl><![CDATA[picurl]]></PicUrl><Url><![CDATA[url]]></Url></item>",
		},
		{
			name: "2",
			args: args{
				title:  "title2",
				desc:   "description2",
				picURL: "picurl",
				url:    "url",
			},
			want: "<item><Title><![CDATA[title2]]></Title><Description><![CDATA[description2]]></Description><PicUrl><![CDATA[picurl]]></PicUrl><Url><![CDATA[url]]></Url></item>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNewsItemMsg(tt.args.title, tt.args.desc, tt.args.picURL, tt.args.url); got != tt.want {
				t.Errorf("NewNewsItemMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNewsReplyMsg(t *testing.T) {
	type args struct {
		toUserID   string
		fromCorpID string
		items      []string
		creatTime  int64
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
				toUserID:   "toUser",
				fromCorpID: "fromUser",
				items: []string{
					"<item><Title><![CDATA[title1]]></Title><Description><![CDATA[description1]]></Description><PicUrl><![CDATA[picurl]]></PicUrl><Url><![CDATA[url]]></Url></item>",
					"<item><Title><![CDATA[title]]></Title><Description><![CDATA[description]]></Description><PicUrl><![CDATA[picurl]]></PicUrl><Url><![CDATA[url]]></Url></item>",
				},
				creatTime: 12345678,
			},
			want: "<xml><ToUserName><![CDATA[toUser]]></ToUserName><FromUserName><![CDATA[fromUser]]></FromUserName><CreateTime>12345678</CreateTime><MsgType><![CDATA[news]]></MsgType><ArticleCount>2</ArticleCount><Articles><item><Title><![CDATA[title1]]></Title><Description><![CDATA[description1]]></Description><PicUrl><![CDATA[picurl]]></PicUrl><Url><![CDATA[url]]></Url></item><item><Title><![CDATA[title]]></Title><Description><![CDATA[description]]></Description><PicUrl><![CDATA[picurl]]></PicUrl><Url><![CDATA[url]]></Url></item></Articles></xml>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNewsReplyMsg(tt.args.toUserID, tt.args.fromCorpID, tt.args.items, tt.args.creatTime); got != tt.want {
				t.Errorf("NewNewsReplyMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
