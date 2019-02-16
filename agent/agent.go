package agent

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qingtao/wxcorp/corp"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

const (
	// 重试间隔:1秒
	retryInterval     = 1 * time.Second
	retryTimes        = 2
	refreshBefore     = 5 * 60 // 秒
	letterForNonceStr = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	minLength         = 16
)

// generateTimestampAndNonceStr 生成时间戳和随机字符
func generateTimestampAndNonceStr(n int) (timestamp int64, s string) {
	source := []byte(letterForNonceStr)
	length := len(source)
	if n < minLength {
		n = minLength
	} else if n > length {
		n = length
	}
	timestamp = time.Now().Unix()
	rand.Seed(timestamp)
	rand.Shuffle(length, func(i, j int) {
		source[i], source[j] = source[j], source[i]
	})
	s = string(source[:n])
	return
}

// Agent 企业微信应用
type Agent struct {
	sync.Mutex
	// AgentID 应用ID
	AgentID string `json:"agentid"`
	// CorpID 企业ID
	CorpID string `json:"cropid"`
	// Secret 应用秘钥
	Secret string `json:"secret"`
	// EncodingAESKey AES秘钥
	EncodingAESKey string `json:"encoding_aes_key"`
	// Token 预留的校验码
	Token string `json:"token"`

	// ipList 微信企业号服务器的ip地址白名单
	ipList map[string]struct{}

	// accessToken 微信企业好应用的访问令牌
	accessToken accessToken
	// refreshAccessToken 刷新令牌标记
	refreshAccessToken int32

	// jsAPITicket jsapi_ticket
	jsAPITicket jsAPITicket
	// 刷新jsapi_ticket标记
	refreshJsAPITicket int32

	// agentJsAPITicket
	agentJsAPITicket jsAPITicket
	// 刷新agent jsapi_ticket标记
	refreshAgentJsAPITicket int32
}

// AccessToken access_token
type accessToken struct {
	accessToken string
	expiresAt   int64
}

// JsAPITiket jspaiticket
type jsAPITicket struct {
	ticket    string
	expiresAt int64
}

// NewAgent 新建企业号APP
func NewAgent(corpid, agentid, secret, encodingAESKey, token string) *Agent {
	return &Agent{
		CorpID:         corpid,
		AgentID:        agentid,
		Secret:         secret,
		EncodingAESKey: encodingAESKey,
		Token:          token,
		ipList:         make(map[string]struct{}),
	}
}

// GetAccessToken 读取AccessToken
func (a *Agent) GetAccessToken() (token string, err error) {
	accessToken, ct := a.accessToken, time.Now().Unix()
	if accessToken.accessToken != "" && ct < accessToken.expiresAt {
		token = accessToken.accessToken
		return
	}
	return a.RefreshAccessToken()
}

// var aonce sync.Once

// RefreshAccessToken 刷新访问令牌
func (a *Agent) RefreshAccessToken() (token string, err error) {
	// println("--刷新访问令牌--", a.Token)
	if !atomic.CompareAndSwapInt32(&a.refreshAccessToken, 0, 1) {
		token = a.accessToken.accessToken
		return
	}
	defer atomic.StoreInt32(&a.refreshAccessToken, 0)
	accesstoken, err := corp.GetAccessToken("", a.CorpID, a.Secret)
	if err != nil {
		return "", err
	}

	token = accesstoken.AccessToken
	// 提前5分钟刷新
	if accesstoken.ExpiresIn > refreshBefore {
		accesstoken.ExpiresIn -= refreshBefore
	}
	a.Lock()
	a.accessToken = accessToken{accesstoken.AccessToken, time.Now().Add(time.Second * time.Duration(accesstoken.ExpiresIn)).Unix()}
	a.Unlock()
	return
}

// GetJsAPITicket 读取jsapi_ticket
func (a *Agent) GetJsAPITicket(typ string) (ticket string, err error) {
	var jsTicket jsAPITicket
	switch typ {
	case "agent_config":
		jsTicket = a.agentJsAPITicket
	default:
		jsTicket = a.jsAPITicket
	}
	ct := time.Now().Unix()
	if jsTicket.ticket != "" && ct < jsTicket.expiresAt {
		ticket = jsTicket.ticket
		return
	}
	return a.RefreshJsAPITicket(typ)
}

// RefreshJsAPITicket -
func (a *Agent) RefreshJsAPITicket(typ string) (ticket string, err error) {
	// println("--刷新jsapi_ticket--", a.Token)
	flag := &a.refreshJsAPITicket
	if typ == "agent_config" {
		flag = &a.refreshAgentJsAPITicket
	}
	if !atomic.CompareAndSwapInt32(flag, 0, 1) {
		if typ == "agent_config" {
			ticket = a.agentJsAPITicket.ticket
		} else {
			ticket = a.jsAPITicket.ticket
		}
		return
	}
	defer atomic.StoreInt32(flag, 0)
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return "", err
	}
	jsTicket, err := corp.GetJsAPITicket("", accessToken, typ)
	if err != nil {
		return "", err
	}
	// 提前5分钟刷新
	if jsTicket.ExpiresIn > refreshBefore {
		jsTicket.ExpiresIn -= refreshBefore
	}
	ticket = jsTicket.Ticket
	jsAPITicket := jsAPITicket{ticket, time.Now().Add(time.Second * time.Duration(jsTicket.ExpiresIn)).Unix()}
	a.Lock()
	if typ == "agent_config" {
		a.agentJsAPITicket = jsAPITicket
	} else {
		a.jsAPITicket = jsAPITicket
	}
	a.Unlock()
	return
}

// NewJsAPITicketSignature 生成ticket签名
func (a *Agent) NewJsAPITicketSignature(ticket, url string) *corp.JsAPITicketSignature {
	timestamp, noncestr := generateTimestampAndNonceStr(0)
	return corp.NewJsAPITicketSignature(a.CorpID, a.AgentID, ticket, noncestr, url, timestamp)
}

// SetIPList 设置IP白名单
func (a *Agent) SetIPList() error {
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return err
	}
	ipList, err := corp.GetCallBackIPList("", accessToken)
	if err != nil {
		return err
	}
	a.Lock()
	defer a.Unlock()
	for _, ip := range ipList.IPList {
		a.ipList[ip] = struct{}{}
	}
	return nil
}

// IsMatch 匹配ip列表
func (a *Agent) IsMatch(remoteAddr string) bool {
	// 如果ip列表为空,直接返回真以便不影响服务
	if len(a.ipList) < 1 {
		return true
	}
	index := strings.LastIndex(remoteAddr, ".")
	ip := remoteAddr[:index] + ".*"
	if _, ok := a.ipList[ip]; ok {
		return true
	}
	return false

}

// ReceiveMsg 接收消息
func (a *Agent) ReceiveMsg(r *http.Request) (msg []byte, err error) {
	crypt := wxbizmsgcrypt.NewWXBizMsgCrypt(a.Token, a.EncodingAESKey, a.CorpID, wxbizmsgcrypt.XmlType)
	signature, nonce, timestamp := r.FormValue("msg_signature"), r.FormValue("nonce"), r.FormValue("timestamp")
	if r.Body != nil {
		defer r.Body.Close()
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var cryptErr *wxbizmsgcrypt.CryptError
	msg, cryptErr = crypt.DecryptMsg(signature, timestamp, nonce, b)
	if cryptErr != nil {
		err = fmt.Errorf("[errcode]:%d,[errmsg]:%s", cryptErr.ErrCode, cryptErr.ErrMsg)
	}
	return
}
