package corp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/qingtao/wxcorp/corp/errcode"
)

const (
	defaultGetUserURL       = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	defaultGetSimpleListURL = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist"
	defaultGetUserListURL   = "https://qyapi.weixin.qq.com/cgi-bin/user/list"
)

// UserResponse 请求用户的响应结构
type UserResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// User 用户信息
	User
}

// Validate 验证响应
func (res *UserResponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// ExtAttr 扩展属性
type ExtAttr struct {
	Type        int            `json:"type,omitempty" xml:",omitempty"`
	Name        string         `json:"name,omitempty" xml:",omitempty"`
	Text        ExtText        `json:"text,omitempty" xml:",omitempty"`
	Web         ExtWeb         `json:"web,omitempty" xml:",omitempty"`
	Miniprogram ExtMiniprogram `json:"miniprogram,omitempty"`
}

// ExtText 扩展属性文本
type ExtText struct {
	Value string `json:"value,omitempty" xml:",omitempty"`
}

// ExtWeb 扩展属性的网页格式
type ExtWeb struct {
	Title string `json:"title,omitempty" xml:",omitempty"`
	URL   string `json:"url,omitempty" xml:",omitempty"`
}

// ExtMiniprogram 扩展属性的小程序格式
type ExtMiniprogram struct {
	Title    string `json:"title,omitempty" xml:",omitempty"`
	AppID    string `json:"appid,omitempty" xml:",omitempty"`
	PagePath string `json:"pagepath,omitempty" xml:",omitempty"`
}

// ExtAttrs 扩展字段
type ExtAttrs struct {
	Attrs []ExtAttr `json:"attrs,omitempty" xml:",omitempty"`
}

// User 用户信息
type User struct {
	// UserID 	成员UserID,对应管理端的帐号,企业内部必须唯一，1-64个字节
	UserID    string `json:"userid"`
	NewUserID string `json:"-" xml:",omitempty"`
	// Name 成员名称,1-64个utf8字符
	Name string `json:"name"`
	// Department 成员所属部门id列表,不超过20个
	Department []int `json:"department"`
	// IsLeaderInDept 是否是部门领导，与Department字段一一对应
	IsLeaderInDept []int `json:"is_leader_in_dept,omitempty"`
	// Position	职务信息,0-128个字符
	Position string `json:"position,omitempty"`
	// Mobile 手机号码,第三方仅通讯录套件可获取,企业内必须唯一，mobile/email二者不能同时为空
	Mobile string `json:"mobile,omitempty"`
	// Gender 性别 0表示未定义，1表示男性，2表示女性
	Gender string `json:"gender,omitempty"`
	// Email 邮箱 第三方仅通讯录套件可获取,6-64个字符
	Email string `json:"email,omitempty"`
	// Avatar 头像url 注：如果要获取小图将url最后的"/0"改成"/64"即可
	Avatar string `json:"avatar,omitempty" xml:",omitempty"`
	// AvatarMediaID 头像ID
	AvatarMediaID string `json:"avatar_mediaid,omitempty" xml:",omitempty"`
	// Telphone 电话号码,32字节以内，可以包含数字和"-"
	Telephone string `json:"telephone,omitempty"`
	// Alias 别名
	Alias string `json:"alias,omitempty"`
	// Status 关注状态: 1=已关注，2=已禁用，4=未关注
	Status int `json:"status,omitempty"`
	// ExtAttr 扩展属性 第三方仅通讯录套件可获取
	ExtAttr *ExtAttrs `json:"extattr,omitempty" xml:">Item"`
	// ToInvite 是否邀请该成员使用企业微信
	ToInvite bool `json:"to_invite,omitempty" xml:",omitempty"`
	// ExternalPosition 成员对外职务,长度最大12个汉字
	ExternalPosition string `json:"external_position,omitempty" xml:",omitempty"`
	// ExternalProfile 成员对外属性
	ExternalProfile *ExternalProfile `json:"external_profile,omitempty" xml:",omitempty"`
}

// ExternalProfile -
type ExternalProfile struct {
	//
	ExternalCoprName string    `json:"external_corp_name,omitempty" xml:",omitempty"`
	ExternalAttr     []ExtAttr `json:"external_attr,omitempty" xml:",omitempty"`
}

// NewGetUserURL 新建获取成员的URL
func NewGetUserURL(url, accessToken, userid string) string {
	if accessToken == "" || userid == "" {
		return ""
	}
	if url == "" {
		url = defaultGetUserURL
	}
	return fmt.Sprintf("%s?access_token=%s&userid=%s", url, accessToken, userid)
}

// GetUser 获取user信息
func GetUser(url, accessToken string, userid string) (user *UserResponse, err error) {
	url = NewGetUserURL(url, accessToken, userid)
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
	if err = json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	err = user.Validate()
	return
}

func sliceIntRemoveDuplicate(a []int) (b []int) {
	if len(a) < 2 {
		return a
	}
	exists := make(map[int]struct{})
	for _, v := range a {
		if _, ok := exists[v]; ok {
			continue
		}
		exists[v] = struct{}{}
		b = append(b, v)
	}
	return
}

// NewGetUserListURL 新建请求部门成员的URL, url主要用户测试，外部调用一般留空即可
func NewGetUserListURL(url, typ, accessToken string, departmentID, fetchChild int, status []int) string {
	if accessToken == "" || departmentID < 1 {
		return ""
	}
	if url == "" {
		if typ == "simple" {
			url = defaultGetSimpleListURL
		} else { // 如果未提供请求路径(不包含请求参数),默认请求用户详情
			url = defaultGetUserListURL
		}
	}
	url = fmt.Sprintf("%s?access_token=%s&department_id=%d", url, accessToken, departmentID)
	if fetchChild == 1 {
		url += fmt.Sprintf("&fetch_child=%d", fetchChild)
	}
	if len(status) < 1 {
		return url
	}
	status = sliceIntRemoveDuplicate(status)
	for _, st := range status {
		url += fmt.Sprintf("&status=%d", st)
	}
	return url
}

// UserListResponse 请求部门用户列表的响应结构
type UserListResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	UserList []User `json:"userlist"`
}

// Validate 检查响应
func (res *UserListResponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// GetUserList 获取部门用户, 如果typ="simple"查询部门成员, 如果typ!="simple"查询部门成员详情
func GetUserList(url, typ, accessToken string, departmenID, fetchChild int, status []int) (res *UserListResponse, err error) {
	url = NewGetUserListURL(url, typ, accessToken, departmenID, fetchChild, status)
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
