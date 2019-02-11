package corp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/qingtao/wxcorp/corp/errcode"
)

// defaultGetCallBackIPURL 默认的获取企业微信服务器IP的url(不包含参数)
const defaultGetCallBackIPURL = "https://qyapi.weixin.qq.com/cgi-bin/getcallbackip"

// IPListResponse 响应
type IPListResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// IPList IP地址
	IPList []string `json:"ip_list"`
}

// Validate 验证响应
func (ipList *IPListResponse) Validate() error {
	if ipList == nil {
		return ErrIsNil
	}
	return errcode.Error(ipList.ErrCode)
}

// NewIPListURL IP地址段的请求URL
func NewIPListURL(url, accessToken string) (s string) {
	if accessToken == "" {
		return
	}
	if url == "" {
		url = defaultGetCallBackIPURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// GetCallBackIPList 获取企业微信服务器地址段
func GetCallBackIPList(url, accessToken string) (ipList *IPListResponse, err error) {
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
	url = NewIPListURL(url, accessToken)
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
	if err = json.Unmarshal(b, &ipList); err != nil {
		return nil, err
	}
	err = ipList.Validate()
	return
}
