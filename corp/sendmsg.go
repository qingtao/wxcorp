package corp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/qingtao/wxcorp/corp/errcode"
)

const (
	defaultSendMsgURL              = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	mimeApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"
)

// Msg 发送的消息结构体
type Msg struct {
	// ToUser 目标用户的UserId,多个用户使用"|"分割,如果是"@all"标记，则表示全部用户,每次最多1000
	ToUser string `json:"touser,omitempty"`
	// ToParty 目标部门的ID,多个部门使用"|"分割,如果ToUser为"@all",忽略此字段,每次最多100
	ToParty string `json:"toparty,omitempty"`
	// ToTag 目标的标签, 多个标签使用"|"分割,如果ToUser为"@all",忽略此字段, 每次最多100
	ToTag string `json:"totag,omitempty"`
	// MsgType 消息类型, 支持"text|image|video|voice|file|news|mpnews|textcard|markdown"等的一个
	MsgType string `json:"msgtype"`
	// AgentID app的id
	AgentID int `json:"agentid"`
	// 是否是保密消息 0：否,1:是默认0
	Safe int `json:"safe,omitempty"`

	// 消息结构
	Text  *TextMsg  `json:"text,omitempty"`
	Image *MediaMsg `json:"image,omitempty"`
	Voice *MediaMsg `json:"voice,omitempty"`
	Video *MediaMsg `json:"video,omitempty"`
	File  *MediaMsg `json:"file,omitempty"`
	// 文本卡片
	TextCard *TextCardMsg `json:"textcard,omitempty"`
	// 图文
	News   *NewsMsg   `json:"news,omitempty"`
	MpNews *MpNewsMsg `json:"mpnews,omitempty"`
	// Markdown 的子集
	Markdown *MarkdownMsg `json:"markdown,omitempty"`
	// 小程序消息暂不支持
}

// Validate 验证消息
func (msg *Msg) Validate() error {
	if msg.ToUser == "" && msg.ToParty == "" && msg.ToTag == "" {
		return errors.New("touser/toparty/totag不能同时为空")
	}
	if msg.AgentID == 0 {
		return errors.New("应用代理agentid不可为空")
	}
	if msg == nil {
		return errors.New("消息为空")
	}
	switch msg.MsgType {
	case "text":
		return msg.Text.Validate()
	case "image":
		return msg.Image.Validate()
	case "voice":
		return msg.Voice.Validate()
	case "video":
		return msg.Video.Validate()
	case "file":
		return msg.File.Validate()
	case "news":
		return msg.News.Validate()
	case "mpnews":
		return msg.MpNews.Validate()
	case "textcard":
		return msg.TextCard.Validate()
	case "markdown":
		return msg.Markdown.Validate()
	}
	return errors.New("消息类型错误")
}

// TextMsg 文本消息
type TextMsg struct {
	Content string `json:"content"`
}

// Validate 验证文本消息
func (msg *TextMsg) Validate() error {
	if msg == nil {
		return errors.New("消息为空")
	}
	if msg.Content == "" {
		return errors.New("文本消息内容为空")
	}
	if len(msg.Content) > 2048 {
		return errors.New("文本消息内容长度超过2048个字节")
	}
	return nil
}

// MediaMsg 图片|语音|视频|文件消息
type MediaMsg struct {
	MediaID string `json:"media_id,omitempty"`
	Title   string `json:"title,omitempty"`
	Desc    string `json:"description,omitempty"`
}

// Validate 验证多媒体消息
func (msg *MediaMsg) Validate() error {
	if msg == nil {
		return errors.New("消息为空")
	}
	if msg.MediaID == "" {
		return errors.New("消息MediaID为空")
	}
	return nil

}

// TextCardMsg 文本卡片消息
type TextCardMsg struct {
	Title  string `json:"title,omitempty"`
	Desc   string `json:"description,omitempty"`
	URL    string `json:"url,omitempty"`
	Btntxt string `json:"btntxt,omitempty"`
}

// Validate 验证卡片消息
func (msg *TextCardMsg) Validate() error {
	if msg == nil {
		return errors.New("消息为空")
	}
	if msg.Title == "" {
		return errors.New("卡片消息标题为空")
	} else if len(msg.Title) > 128 {
		return errors.New("卡片消息标题长度超过128字节")
	}
	if msg.Desc == "" {
		return errors.New("卡片消息的描述为空")
	} else if len(msg.Desc) > 512 {
		return errors.New("卡片消息描述长度超过512字节")
	}
	if msg.URL == "" {
		return errors.New("卡片消息链接为空")
	}
	// if len(msg.Btntxt) > 16 {
	// 	return errors.New("卡片消息按钮长度超过限制4个文字")
	// }
	return nil
}

// NewsMsg 图文消息
type NewsMsg struct {
	Acticles []NewsItem `json:"acticles"`
}

// Validate 验证图文消息
func (msg *NewsMsg) Validate() error {
	if msg == nil {
		return errors.New("图文消息为空")
	}
	if len(msg.Acticles) < 1 || len(msg.Acticles) > 8 {
		return errors.New("图文消息支持1到8条图文")
	}
	var err error
	for _, article := range msg.Acticles {
		if err = article.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// NewsItem 图文消息的项目
type NewsItem struct {
	Title  string `json:"title,omitempty"`
	Desc   string `json:"description,omitempty"`
	URL    string `json:"url,omitempty"`
	PicURL string `json:"picurl,omitempty"`
}

// Validate 验证图文消息
func (item NewsItem) Validate() error {
	if item.Title == "" {
		return errors.New("图文消息标题为空")
	} else if len(item.Title) > 128 {
		return errors.New("图文消息标题长度超过128字节")
	}
	if item.URL == "" {
		return errors.New("图文消息的链接为空")
	}
	return nil
}

// MpNewsMsg 图文消息
type MpNewsMsg struct {
	Acticles []MpNewsItem `json:"acticles"`
}

// Validate 验证mpnews
func (msg *MpNewsMsg) Validate() error {
	if msg == nil {
		return errors.New("图文消息为空")
	}
	if len(msg.Acticles) < 1 || len(msg.Acticles) > 8 {
		return errors.New("图文消息支持1到8条图文")
	}
	var err error
	for _, article := range msg.Acticles {
		if err = article.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// MpNewsItem mpnews消息与news消息类似，不同的是图文消息内容存储在微信后台，并且支持保密选项。每个应用每天最多可以发送100次
type MpNewsItem struct {
	Title            string `json:"title,omitempty"`
	ThumbMediaID     string `json:"thumb_media_id,omitempty"`
	Author           string `json:"author,omitempty"`
	ContentSourceURL string `json:"content_source_url,omitempty"`
	Content          string `json:"content,omitempty"`
	Digest           string `json:"digest,omitempty"`
	ShowCoverPic     string `json:"show_cover_pic,omitempty"`
}

// Validate 验证图文消息
func (item MpNewsItem) Validate() error {
	if item.Title == "" {
		return errors.New("图文消息标题为空")
	} else if len(item.Title) > 128 {
		return errors.New("图文消息标题长度超过128字节")
	}
	if item.ThumbMediaID == "" {
		return errors.New("图文消息的链接为空")
	}
	if item.Content == "" {
		return errors.New("图文消息内容为空")
	} else if len(item.Content) > 666 {
		return errors.New("图文消息内容长度超过666个字节")
	}
	if len(item.Author) > 64 {
		return errors.New("图文消息作者超过64个字节")
	}
	if len(item.Digest) > 512 {
		return errors.New("图文消息描述超过512个字节")
	}
	return nil
}

// MarkdownMsg markdown格式的消息
// 只支持的语法:
//	1. 1-6级标题 #
//	2. 加粗 **
//	3. 链接 [link](url)
//	4. 行内代码段(不支持换行) `code`
//	5. 引用 >
// 只支持的文字颜色:
//	1. <font color="info">绿色</font>
//	2. <font color="comment">灰色</font>
//	3. <font color="warning">橙红色</font>
type MarkdownMsg struct {
	Content string `json:"content,omitempty"`
}

// Validate 验证markdown消息
func (msg *MarkdownMsg) Validate() error {
	if msg == nil {
		return errors.New("消息为空")
	}
	if msg.Content == "" {
		return errors.New("markdown消息内容为空")
	}
	if len(msg.Content) > 2048 {
		return errors.New("markdown消息内容长度超过2048个字节")
	}
	return nil
}

// SendMsgResponse 发送消息的响应结构
//	收件人必须处于应用的可见范围内，并且管理组对应用有使用权限、对收件人有查看权限，否则本次调用失败。
//	如果无权限或收件人不存在，则本次发送失败，返回无效的userid列表（注：由于userid不区分大小写，返回的列表都统一转为小写）；如果未关注，发送仍然执行。
type SendMsgResponse struct {
	ErrCode      int    `json:"errcode,omitempty"`
	ErrMsg       string `json:"errmsg,omitempty"`
	InvalidUser  string `json:"invaliduser,omitempty"`
	InvalidParty string `json:"invalidparty,omitempty"`
	InvalidTag   string `json:"ivalidtag,omitempty"`
}

// NewSendMsgURL 新建发送消息URL
func NewSendMsgURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultSendMsgURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// SendMsg 发送应用消息
func SendMsg(url, accessToken string, msg *Msg) (err error) {
	if accessToken == "" {
		return errcode.ErrInvalidAccessToken
	}
	if err = msg.Validate(); err != nil {
		return
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	buf := bytes.NewBuffer(b)
	defer buf.Reset()
	url = NewSendMsgURL(url, accessToken)
	resp, err := httpClient.Post(url, mimeApplicationJSONCharsetUTF8, buf)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var res SendMsgResponse
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}

	if err = errcode.Error(res.ErrCode); err != nil {
		return
	}
	var errStr string
	if len(res.InvalidUser) > 0 {
		errStr += "[InvalidUser]:" + res.InvalidUser
	}
	if len(res.InvalidParty) > 0 {
		errStr += "[InvalidParty]:" + res.InvalidParty
	}
	if len(res.InvalidTag) > 0 {
		errStr += "[InvalidTag]:" + res.InvalidTag
	}
	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}
