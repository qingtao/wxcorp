package corp

import (
	"fmt"
)

// Request 是微信推送的事件：如关注、取消关注，消息：如文本、图片类型的消息
type Request struct {
	// 微信请求头
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	AgentID      string
	// MsgID 可用于普通消息排重
	MsgID int64 `xml:"MsgId"`

	//消息内容
	Content      string
	PicURL       string `xml:"PicUrl"`
	MediaID      string `xml:"MediaId"`
	Format       string
	ThumbMediaID string  `xml:"ThumbMediaId"`
	LocationX    float32 `xml:"Location_X"`
	LocationY    float32 `xml:"Location_Y"`
	Scale        float32
	Label        string
	Title        string
	Description  string
	URL          string `xml:"Url"`
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float32
	Longitude    float32
	Precision    float32
	Recognition  string
	Status       string
	// scancode info
	ScanCodeInfo ScanCodeInfo
	// image info
	SendPicsInfo SendPicsInfo
	// location info
	SendLocationInfo SendLocationInfo
	// 通讯录事件单独处理
	// ContactEvent `xml:",omitempty"`
}

// BatchJob 异步任务
type BatchJob struct {
	JobID    string `xml:"JobId"`
	JobType  string
	ErroCode int    // 错误码
	ErrMsg   string // 错误消息
}

// ScanCodeInfo button scancode_push or scancode_waitmsg event
type ScanCodeInfo struct {
	// ScanType qrcode
	ScanType   string
	ScanResult string
}

// SendPicsInfo button pic_photo_or_album message, pic_sysphoto or pic_weixin event
type SendPicsInfo struct {
	Count   int
	PicList []PicItem `xml:"PicList>item"`
}

// PicItem image md5sum
type PicItem struct {
	PicMd5Sum string `xml:"PicMd5Sum"`
}

// SendLocationInfo button location_select event
type SendLocationInfo struct {
	LocationX float32 `xml:"Location_X"`
	LocationY float32 `xml:"Location_Y"`
	Scale     float32
	Label     string
	Poiname   string
}

// NewTextReplyMsg 新建文本消息
func NewTextReplyMsg(toUserID, fromCorpID, content string, createTime int64) string {
	return fmt.Sprintf("<xml><ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[%s]]></Content></xml>", toUserID, fromCorpID, createTime, content)
}

// NewImageReplyMsg 新建图片消息
func NewImageReplyMsg(toUserID, fromCorpID, mediaID string, createTime int64) string {
	return fmt.Sprintf("<xml><ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime><MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[%s]]></MediaId></Image></xml>", toUserID, fromCorpID, createTime, mediaID)
}

// NewVoiceReplyMsg 新建语音消息
func NewVoiceReplyMsg(toUserID, fromCorpID, mediaID string, createTime int64) string {
	return fmt.Sprintf("<xml><ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime><MsgType><![CDATA[voice]]></MsgType><Voice><MediaId><![CDATA[%s]]></MediaId></Voice></xml>", toUserID, fromCorpID, createTime, mediaID)
}

// NewVideoReplyMsg 新建视频消息
func NewVideoReplyMsg(toUserID, fromCorpID, mediaID, title, desc string, createTime int64) string {
	return fmt.Sprintf("<xml><ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime><MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[%s]]></MediaId><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description></Video></xml>", toUserID, fromCorpID, createTime, mediaID, title, desc)
}

// NewNewsItemMsg 新建图文消息内容
func NewNewsItemMsg(title, desc, picURL, url string) string {
	return fmt.Sprintf("<item><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description><PicUrl><![CDATA[%s]]></PicUrl><Url><![CDATA[%s]]></Url></item>", title, desc, picURL, url)
}

// NewNewsReplyMsg 新建图文消息, items是每条图文消息的字符串, 即NewNewsItemMsg的结果
func NewNewsReplyMsg(toUserID, fromCorpID string, items []string, creatTime int64) string {
	var s string
	for _, item := range items {
		s += item
	}
	return fmt.Sprintf("<xml><ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime><MsgType><![CDATA[news]]></MsgType><ArticleCount>%d</ArticleCount><Articles>%s</Articles></xml>", toUserID, fromCorpID, creatTime, len(items), s)
}
