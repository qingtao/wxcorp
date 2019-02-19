package corp

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/qingtao/wxcorp/corp/errcode"
)

const (
	// defaultGetJsAPITicketURL 请求jsapi_ticket的URL(不包含参数)
	defaultGetJsAPITicketURL = "https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket"
	// defaultGetAgentJsAPITicketURL 请求应用的jsapi_ticket(不包含参数)
	defaultGetAgentJsAPITicketURL = "https://qyapi.weixin.qq.com/cgi-bin/ticket/get"
)

// JsAPITicketResponse 企业的jsapi_ticket响应结果
type JsAPITicketResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// Ticket 生成签名所需的jsapi_ticket
	Ticket string `json:"ticket"`
	// ExpiresIn 有效时间,单位（秒）
	ExpiresIn int `json:"expires_in"`
}

// Validate 验证响应
func (res *JsAPITicketResponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// NewJsAPITicketURL 新建请求jsapi_ticket的URL
func NewJsAPITicketURL(url, accessToken, typ string) string {
	if accessToken == "" {
		return ""
	}
	if typ == "agent_config" {
		if url == "" {
			url = defaultGetAgentJsAPITicketURL
		}
		return fmt.Sprintf("%s?access_token=%s&type=agent_config", url, accessToken)
	}
	if url == "" {
		url = defaultGetJsAPITicketURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// GetJsAPITicket 请求jsapi_ticket
func GetJsAPITicket(url, accessToken, typ string) (ticket *JsAPITicketResponse, err error) {
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
	url = NewJsAPITicketURL(url, accessToken, typ)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &ticket); err != nil {
		return
	}
	err = ticket.Validate()
	return
}

// genSignature 生成ticket的签名
func genSignature(ticket, noncestr, url string, timestamp int64) string {
	s := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, noncestr, timestamp, url)
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

// JsAPITicketSignature jsapi_ticket的签名
type JsAPITicketSignature struct {
	CorpID    string `json:"corpid"`
	AgentID   string `json:"agentid,omitempty"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

// NewJsAPITicketSignature 新建jsapi_ticket签名
func NewJsAPITicketSignature(corpid, agentid, ticket, noncestr, url string, timestamp int64) *JsAPITicketSignature {
	return &JsAPITicketSignature{
		CorpID:    corpid,
		AgentID:   agentid,
		Timestamp: timestamp,
		NonceStr:  noncestr,
		Signature: genSignature(ticket, noncestr, url, timestamp),
	}
}
