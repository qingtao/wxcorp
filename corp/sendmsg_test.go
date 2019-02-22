package corp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

func TestMsg_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *Msg
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			msg: &Msg{
				ToUser: "",
			},
			wantErr: true,
		},
		{
			name: "2",
			msg: &Msg{
				ToUser:  "1",
				AgentID: 0,
			},
			wantErr: true,
		},
		{
			name:    "3",
			msg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.msg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Msg.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	msg := &Msg{
		ToUser:  "TOUSER",
		AgentID: 1002,
		MsgType: "text",
		Text: &TextMsg{
			Content: "abc",
		},
	}
	media := &MediaMsg{
		MediaID: "1",
	}
	t.Run("text", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})
	msg.MsgType = "image"
	msg.Image = media
	t.Run("image", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})
	msg.MsgType = "voice"
	msg.Voice = media
	t.Run("voice", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})
	msg.MsgType = "video"
	msg.Video = media
	t.Run("video", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})
	msg.MsgType = "file"
	msg.File = media
	t.Run("file", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})
	msg.MsgType = "news"
	var s = `{"articles":[{"title":"中秋节礼品领取","description":"今年中秋节公司有豪礼相送","url":"URL","picurl":"http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png"}]}`
	var news NewsMsg
	json.Unmarshal([]byte(s), &news)
	msg.News = &news
	t.Run("news", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})

	msg.MsgType = "mpnews"
	s = `{"articles":[{"title":"Title","thumb_media_id":"MEDIA_ID","author":"Author","content_source_url":"URL","content":"Content","digest":"Digest description"}]}`
	var mpnews MpNewsMsg
	json.Unmarshal([]byte(s), &mpnews)
	msg.MpNews = &mpnews
	t.Run("mpnews", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})

	msg.MsgType = "textcard"
	s = `{"title":"领奖通知","description":"<div class=\"gray\">2016年9月26日</div> <div class=\"normal\">恭喜你抽中iPhone 7一台，领奖码：xxxx</div><div class=\"highlight\">请于2016年10月10日前联系行政同事领取</div>","url":"URL","btntxt":"更多"}`
	var textCardMsg TextCardMsg
	json.Unmarshal([]byte(s), &textCardMsg)
	msg.TextCard = &textCardMsg
	t.Run("textcard", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})

	s = `{
		"content": "您的会议室已经预定，稍后会同步到邮箱                 >**事项详情**                  >事　项：<font color=\"info\">开会</font>                  >组织者：@miglioguan                  >参与者：@miglioguan、@kunliu、@jamdeezhou、@kanexiong、@kisonwang                  >                  >会议室：<font color=\"info\">广州TIT 1楼 301</font>                  >日　期：<font color=\"warning\">2018年5月18日</font>                  >时　间：<font color=\"comment\">上午9:00-11:00</font>                  >                  >请准时参加会议。                  >                  >如需修改会议信息，请点击：[修改会议信息](https://work.weixin.qq.com)"
	}`
	var mdMsg MarkdownMsg
	msg.MsgType = "markdown"
	json.Unmarshal([]byte(s), &mdMsg)
	msg.Markdown = &mdMsg
	t.Run("markdown", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Error(err.Error())
		}
	})
	msg.MsgType = "invalidType"
	t.Run("invalidType", func(t *testing.T) {
		if err := msg.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("此处应该返回错误")
		}
	})
}

func TestTextMsg_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *TextMsg
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			msg: &TextMsg{
				Content: "",
			},
			wantErr: true,
		},
		{
			name:    "2",
			msg:     nil,
			wantErr: true,
		},
		{
			name: "3",
			msg: &TextMsg{
				Content: strings.Repeat("1234567890", 205),
			},
			wantErr: true,
		},
		{
			name: "4",
			msg: &TextMsg{
				Content: "abc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.msg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("TextMsg.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMediaMsg_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *MediaMsg
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			msg: &MediaMsg{
				MediaID: "1",
				Title:   "标题",
				Desc:    "",
			},
			wantErr: false,
		},
		{
			name: "2",
			msg: &MediaMsg{
				MediaID: "",
			},
			wantErr: true,
		},
		{
			name:    "3",
			msg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.msg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("MediaMsg.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTextCardMsg_Validate(t *testing.T) {
	var s = `{"title":"领奖通知","description":"<div class=\"gray\">2016年9月26日</div> <div class=\"normal\">恭喜你抽中iPhone 7一台，领奖码：xxxx</div><div class=\"highlight\">请于2016年10月10日前联系行政同事领取</div>","url":"URL","btntxt":"更多"}`
	var msg TextCardMsg
	json.Unmarshal([]byte(s), &msg)

	t.Run("empty", func(t *testing.T) {
		var m *TextCardMsg
		if err := m.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("此处应该返回错误")
		}
	})
	title := msg.Title
	t.Run("title", func(t *testing.T) {
		msg.Title = ""
		if err := msg.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("此处应该返回错误")
		}
		t.Run("1", func(t *testing.T) {
			msg.Title = strings.Repeat(title, 12)
			if err := msg.Validate(); err != nil {
				t.Log(err.Error())
			} else {
				t.Error("此处应该返回错误")
			}
		})
	})
	msg.Title = title
	desc := msg.Desc
	t.Run("desc", func(t *testing.T) {
		msg.Desc = ""
		if err := msg.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("此处应该返回错误")
		}
		t.Run("1", func(t *testing.T) {
			msg.Desc = strings.Repeat(desc, 5)
			if err := msg.Validate(); err != nil {
				t.Log(err.Error())
			} else {
				t.Error("此处应该返回错误")
			}
		})
	})
	msg.Desc = desc
	t.Run("url", func(t *testing.T) {
		msg.URL = ""
		if err := msg.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("此处应该返回错误")
		}
	})
}

func TestNewsMsg_Validate(t *testing.T) {
	var s = `{
               "title" : "中秋节礼品领取",
               "description" : "今年中秋节公司有豪礼相送",
               "url" : "URL",
               "picurl" : "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png"
		   }`
	var item NewsItem
	json.Unmarshal([]byte(s), &item)
	articles := &NewsMsg{
		Articles: []NewsItem{item},
	}
	tests := []struct {
		name     string
		articles *NewsMsg
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "1",
			articles: articles,
			wantErr:  false,
		},
		{
			name:     "2",
			articles: nil,
			wantErr:  true,
		},
		{
			name:     "3",
			articles: new(NewsMsg),
			wantErr:  true,
		},
		{
			name: "4",
			articles: &NewsMsg{
				Articles: []NewsItem{
					NewsItem{Title: ""},
				},
			},
			wantErr: true,
		},
		{
			name: "5",
			articles: &NewsMsg{
				Articles: []NewsItem{
					NewsItem{
						Title: "abc",
						URL:   "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "6",
			articles: &NewsMsg{
				Articles: []NewsItem{
					NewsItem{
						Title: strings.Repeat(item.Title, 7),
						URL:   "http://a.b.com",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.articles.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("NewsMsg.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsItem_Validate(t *testing.T) {
	type fields struct {
		Title  string
		Desc   string
		URL    string
		PicURL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewsItem{
				Title:  tt.fields.Title,
				Desc:   tt.fields.Desc,
				URL:    tt.fields.URL,
				PicURL: tt.fields.PicURL,
			}
			if err := item.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("NewsItem.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMpNewsMsg_Validate(t *testing.T) {
	var s = `{"articles":[{"title":"Title","thumb_media_id":"MEDIA_ID","author":"Author","content_source_url":"URL","content":"Content","digest":"Digest description"}]}`
	var articles MpNewsMsg
	json.Unmarshal([]byte(s), &articles)
	tests := []struct {
		name     string
		articles *MpNewsMsg
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "1",
			articles: &articles,
			wantErr:  false,
		},
		{
			name:     "2",
			articles: new(MpNewsMsg),
			wantErr:  true,
		},
		{
			name:     "3",
			articles: nil,
			wantErr:  true,
		},
		{
			name: "4",
			articles: &MpNewsMsg{
				[]MpNewsItem{
					MpNewsItem{Title: ""},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.articles.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("MpNewsMsg.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMpNewsItem_Validate(t *testing.T) {
	var s = `{"title":"Title","thumb_media_id":"MEDIA_ID","author":"Author","content_source_url":"URL","content":"Content","digest":"Digest description"}`
	var item MpNewsItem
	json.Unmarshal([]byte(s), &item)
	t.Run("1", func(t *testing.T) {
		if err := item.Validate(); err != nil {
			t.Error(err)
		}
		t.Run("1", func(t *testing.T) {
			title := item.Title
			item.Title = ""
			if err := item.Validate(); err != nil {
				t.Log(err.Error())
			} else {
				t.Error("应该有错误，但是此处返回错误为空")
			}
			t.Run("1", func(t *testing.T) {
				item.Title = strings.Repeat(title, 26)
				if err := item.Validate(); err != nil {
					t.Log(err.Error())
				} else {
					t.Error("应该有错误，但是此处返回错误为空")
				}
			})
			item.Title = title
		})
	})
	t.Run("2", func(t *testing.T) {
		mediaID := item.ThumbMediaID
		item.ThumbMediaID = ""
		if err := item.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("应该有错误，但是此处返回错误为空")
		}
		item.ThumbMediaID = mediaID
	})
	t.Run("3", func(t *testing.T) {
		content := item.Content
		item.Content = ""
		if err := item.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("应该有错误，但是此处返回错误为空")
		}
		t.Run("33", func(t *testing.T) {
			item.Content = strings.Repeat(content, 112)
			if err := item.Validate(); err != nil {
				t.Log(err.Error())
			} else {
				t.Error("应该有错误，但是此处返回错误为空")
			}
		})
		item.Content = content
	})
	t.Run("4", func(t *testing.T) {
		author := item.Author
		item.Author = strings.Repeat(author, 12)
		if err := item.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("应该有错误，但是此处返回错误为空")
		}
		item.Author = author
	})
	t.Run("5", func(t *testing.T) {
		digest := item.Digest
		item.Digest = strings.Repeat(digest, 86)
		if err := item.Validate(); err != nil {
			t.Log(err.Error())
		} else {
			t.Error("应该有错误，但是此处返回错误为空")
		}
	})
}

func TestMarkdownMsg_Validate(t *testing.T) {
	var s = "12345678"

	tests := []struct {
		name    string
		msg     *MarkdownMsg
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			msg: &MarkdownMsg{
				Content: "",
			},
			wantErr: true,
		},
		{
			name: "2",
			msg: &MarkdownMsg{
				Content: strings.Repeat(s, 257),
			},
			wantErr: true,
		},
		{
			name: "3",
			msg: &MarkdownMsg{
				Content: s,
			},
			wantErr: false,
		},
		{
			name:    "4",
			msg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.msg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("MarkdownMsg.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendMsg(t *testing.T) {
	s := `{
   "touser" : "UserID1|UserID2|UserID3",
   "toparty" : "PartyID1|PartyID2",
   "totag" : "TagID1 | TagID2",
   "msgtype" : "text",
   "agentid" : 1,
   "text" : {
       "content" : "你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。"
   },
   "safe":0
}`
	var msg Msg
	json.Unmarshal([]byte(s), &msg)

	var resStr = `{
   "errcode" : 0,
   "errmsg" : "ok",
   "invaliduser" : "",
   "invalidparty" : "",
   "invalidtag":""
 }`

	var resErr = `{
   "errcode" : 0,
   "errmsg" : "ok",
   "invaliduser" : "userid1|userid2",
   "invalidparty" : "partyid1|partyid2",
   "invalidtag":"tagid1|tagid2"
 }`

	ht := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := resStr
		r.ParseForm()
		accesstoken := r.FormValue("access_token")
		switch accesstoken {
		case "wantOk":
		case "wantJSONErr":
			res = `{"errcode:0,"errmsg":"ok"}`
		case "wantErr":
			res = resErr
		default:
			res = `{"errcode":1,"errmsg":"未知错误"}`
		}
		fmt.Fprint(w, res)
	}))

	type args struct {
		url         string
		accessToken string
		msg         *Msg
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
				accessToken: "",
			},
			wantErr: true,
		},
		{
			name: "wantOk",
			args: args{
				url:         ht.URL,
				accessToken: "wantOk",
				msg:         &msg,
			},
			wantErr: false,
		},
		{
			name: "wantJSONErr",
			args: args{
				url:         ht.URL,
				accessToken: "wantJSONErr",
				msg:         &msg,
			},
			wantErr: true,
		},
		{
			name: "wantErr",
			args: args{
				url:         ht.URL,
				accessToken: "wantErr",
				msg:         &msg,
			},
			wantErr: true,
		},
		{
			name: "wantErr2",
			args: args{
				url:         ht.URL,
				accessToken: "wantErr2",
				msg:         &msg,
			},
			wantErr: true,
		},
		{
			name: "networkErr",
			args: args{
				url:         `http://127.0.0.1:8080`,
				accessToken: "a",
				msg:         &msg,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendMsg(tt.args.url, tt.args.accessToken, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				t.Log(err.Error())
			}
		})
	}
}
