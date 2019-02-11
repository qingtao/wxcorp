package corp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/qingtao/wxcorp/corp/errcode"

	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

const (
	defaultGetTokenURL = `https://qyapi.weixin.qq.com/cgi-bin/gettoken`
)

var (
	// ErrEchoStrURLInvalid 微信验证URL格式错误
	ErrEchoStrURLInvalid = errors.New("微信验证URL格式错误")
	// ErrCorpIDOrSecretIsEmpty corpid或者secret为空
	ErrCorpIDOrSecretIsEmpty = errors.New("corpid或者secret为空")
	// ErrIsNil 空指针
	ErrIsNil = errors.New("对象为空")
)

// 设置http请求超时时间为10秒
var httpClient = http.Client{
	Timeout: 10 * time.Second,
}

// NewAccessTokenURL 新建获取access_token的URL
func NewAccessTokenURL(url, corpid, secret string) string {
	if corpid == "" || secret == "" {
		return ""
	}
	if url == "" {
		url = defaultGetTokenURL
	}
	return fmt.Sprintf("%s?corpid=%s&corpsecret=%s", url, corpid, secret)
}

// GetAccessToken 获取access_token, 返回err ！= nil 如果出现错误
func GetAccessToken(url, corpid, secret string) (res *AccessTokenResponse, err error) {
	if corpid == "" || secret == "" {
		return nil, ErrCorpIDOrSecretIsEmpty
	}
	addr := NewAccessTokenURL(url, corpid, secret)
	resp, err := httpClient.Get(addr)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// b, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(b, &res)
	if err != nil {
		return
	}

	err = res.Validate()
	return
}

// AccessTokenResponse 企业微信的access_token响应
type AccessTokenResponse struct {
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// Validate 验证access_token响应结果
func (token *AccessTokenResponse) Validate() error {
	if token == nil {
		return ErrIsNil
	}
	return errcode.Error(token.ErrCode)
}

// GetEchoStr 获取echostr,请求地址:
//	http://api.3dept.com/?msg_signature=ASDFQWEXZCVAQFASDFASDFSS&timestamp=13500001234&nonce=123412323&echostr=ENCRYPT_STR
func GetEchoStr(corpid, token, encodingAESKey, s string) (echostr []byte, err error) {
	// 解析url
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	// 解析参数列表
	values := u.Query()

	msgsign, timestamp, nonce, encechostr := values.Get("msg_signature"), values.Get("timestamp"), values.Get("nonce"), values.Get("echostr")
	// 任一查询字符串为空返回错误
	if msgsign == "" || timestamp == "" || nonce == "" || encechostr == "" {
		return nil, ErrEchoStrURLInvalid
	}
	// fmt.Println("1:", msgsign, "2:", timestamp, "3:", nonce, "4:", encechostr)
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAESKey, corpid, wxbizmsgcrypt.XmlType)
	var cryptErr *wxbizmsgcrypt.CryptError
	echostr, cryptErr = wxcpt.VerifyURL(msgsign, timestamp, nonce, encechostr)
	if cryptErr != nil {
		// fmt.Println(cryptErr)
		return nil, errors.New(cryptErr.ErrMsg)
	}
	return
}
