package corp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/qingtao/wxcorp/corp/errcode"
)

const (
	defaultGetTagListURL   = "https://qyapi.weixin.qq.com/cgi-bin/tag/list"
	defaultGetUserOfTagURL = "https://qyapi.weixin.qq.com/cgi-bin/tag/get"
)

// Tag 企业号通讯录标签
type Tag struct {
	// TagName 标签名称
	TagName string `json:"tagname"`
	// TagID 标签ID
	TagID int `json:"tagid"`
}

// TagListResponse 标签列表返回结构
type TagListResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	TagList []Tag  `json:"taglist,omitempty"`
}

// Validate -
func (res *TagListResponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// UserlistOfTag 标签成员
type UserlistOfTag struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
}

// MemberOfTagReponse 标签成员响应结构
type MemberOfTagReponse struct {
	ErrCode   int             `json:"errcode"`
	ErrMsg    string          `json:"errmsg"`
	TagName   string          `json:"tagname"`
	UserList  []UserlistOfTag `json:"userlist"`
	PartyList []int           `json:"partylist"`
}

// Member 标签成员
type Member struct {
	TagName   string          `json:"tagname"`
	UserList  []UserlistOfTag `json:"userlist"`
	PartyList []int           `json:"partylist"`
}

// Validate 验证标签列表响应
func (res *MemberOfTagReponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// NewGetTagListURL 新建获取标签列表的URL
func NewGetTagListURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultGetTagListURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// GetTagList 获取标签列表
func GetTagList(url, accessToken string) (res *TagListResponse, err error) {
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
	url = NewGetTagListURL(url, accessToken)
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
		return
	}
	err = res.Validate()
	return
}

// NewGetUserOfTagURL 新建获取标签用户的URL
func NewGetUserOfTagURL(url, accessToken string, tagid int) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultGetUserOfTagURL
	}
	return fmt.Sprintf("%s?access_token=%s&tagid=%d", url, accessToken, tagid)
}

// GetMemberOfTag 获取标签成员
func GetMemberOfTag(url, accessToken string, tagid int) (member *MemberOfTagReponse, err error) {
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
	url = NewGetUserOfTagURL(url, accessToken, tagid)
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
	if err = json.Unmarshal(b, &member); err != nil {
		return
	}
	err = member.Validate()
	return
}
