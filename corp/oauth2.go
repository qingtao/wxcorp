package corp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/qingtao/wxcorp/corp/errcode"
)

const (
	defaultOAuth2AuthorizeURL   = "https://open.weixin.qq.com/connect/oauth2/authorize"
	defaultOAuth2GetUserInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
)

// NewOAuth2RedirectURL 新建网页授权跳转链接
func NewOAuth2RedirectURL(wxurl, appid, returnTo, targetURI, state string) string {
	if wxurl == "" {
		wxurl = defaultOAuth2AuthorizeURL
	}
	// 添加返回的地址
	if returnTo != "" {
		targetURI = fmt.Sprintf("%s?return_to=%s", targetURI, url.QueryEscape(returnTo))
	}
	targetURI = url.QueryEscape(targetURI)

	return fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect", wxurl, appid, targetURI, state)
}

// NewGetUserInfoURL 新建请求用户身份(userid)链接
func NewGetUserInfoURL(wxurl, accessToken, code string) string {
	if wxurl == "" {
		wxurl = defaultOAuth2GetUserInfoURL
	}
	return fmt.Sprintf("%s?access_token=%s&code=%s", wxurl, accessToken, code)
}

// UserInfoResponse 用户信息返回结构
type UserInfoResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	UserID   string `json:"UserId"`
	DeviceID string `json:"DeviceId"`
}

// Validate 校验响应数据
func (res *UserInfoResponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// GetUserInfoWithCode 获取用户基本信息
func GetUserInfoWithCode(url, accessToken, code string) (res *UserInfoResponse, err error) {
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
	url = NewGetUserInfoURL(url, accessToken, code)
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
	if err = json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	err = res.Validate()
	return
}
