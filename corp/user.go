package corp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/qingtao/wxcorp/corp/errcode"
)

const (
	defaultGetUserURL         = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	defaultGetSimpleListURL   = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist"
	defaultGetUserListURL     = "https://qyapi.weixin.qq.com/cgi-bin/user/list"
	defaultCreateUserURL      = "https://qyapi.weixin.qq.com/cgi-bin/user/create"
	defaultUpdateUserURL      = "https://qyapi.weixin.qq.com/cgi-bin/user/update"
	defaultDeleteUserURL      = "https://qyapi.weixin.qq.com/cgi-bin/user/delete"
	defaultBatchDeleteUserURL = "https://qyapi.weixin.qq.com/cgi-bin/user/batchdelete"

	defaultUserIDToOpenIDURL = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_openid"
	defaultOpenIDToUserIDURL = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_userid"

	maxBatchDeleteUserCount = 200 // 批量删除时，一次请求最多可以删除200个用户
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

// Validate 验证扩展属性
func (a ExtAttr) Validate() error {
	switch a.Type {
	case 0:
		return a.Text.Validate()
	case 1:
		return a.Web.Validate()
	case 2:
		return a.Miniprogram.Validate()
	default:
		return errors.New("不支持的扩展属性类型")
	}
}

// ExtText 扩展属性文本
type ExtText struct {
	Value string `json:"value,omitempty" xml:",omitempty"`
}

// Validate 验证文本属性
func (a ExtText) Validate() error {
	if a.Value != "" && utf8.RuneCountInString(a.Value) > 12 {
		return errors.New("文本属性的内容超过12个UTF8字符")
	}
	return nil
}

// ExtWeb 扩展属性的网页格式
type ExtWeb struct {
	Title string `json:"title,omitempty" xml:",omitempty"`
	URL   string `json:"url,omitempty" xml:",omitempty"`
}

// Validate 网页类型的属性校验
func (a ExtWeb) Validate() error {
	switch {
	case a.Title == "" && a.URL == "":
		return nil
	case a.Title != "" && a.URL != "":
		if utf8.RuneCountInString(a.Title) > 12 {
			return errors.New("网页的展示标题长度限制12个UTF8字符")
		}
	default:
		return errors.New("url和title字段要么同时为空表示清除该属性，要么同时不为空")
	}
	return nil
}

// ExtMiniprogram 扩展属性的小程序格式
type ExtMiniprogram struct {
	Title    string `json:"title,omitempty" xml:",omitempty"`
	AppID    string `json:"appid,omitempty" xml:",omitempty"`
	PagePath string `json:"pagepath,omitempty" xml:",omitempty"`
}

// Validate 小程序类型的属性校验
func (a ExtMiniprogram) Validate() error {
	switch {
	case a.Title == "" && a.AppID == "":
		return nil
	case a.Title != "" && a.AppID != "":
		if utf8.RuneCountInString(a.Title) > 12 {
			return errors.New("小程序的展示标题,长度限制12个UTF8字符")
		}
	default:
		return errors.New("appid和title字段要么同时为空表示清除改属性，要么同时不为空")
	}
	return nil
}

// ExtAttrs 扩展字段
type ExtAttrs struct {
	Attrs []ExtAttr `json:"attrs,omitempty" xml:",omitempty"`
}

// Validate 验证扩展字段
func (a *ExtAttrs) Validate() (err error) {
	for _, attr := range a.Attrs {
		if err = attr.Validate(); err != nil {
			return
		}
	}
	return nil
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
	// Order 成员在部门中的排序
	Order []int `json:"order"`
	// IsLeaderInDept 是否是部门领导，与Department字段一一对应
	IsLeaderInDept []int `json:"is_leader_in_dept,omitempty"`
	// Position	职务信息,0-128个字符
	Position string `json:"position,omitempty"`
	// Mobile 手机号码,第三方仅通讯录套件可获取,企业内必须唯一，mobile/email二者不能同时为空
	Mobile string `json:"mobile,omitempty"`
	// Gender 性别 0表示未定义，1表示男性，2表示女性
	Gender string `json:"gender,omitempty"`
	// Enable 是否启用, 1:启用 0:禁用
	Enable int `json:"enable"`
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

// Validate 验证用户字段
func (a User) Validate() error {
	if a.Enable != 0 && a.Enable != 1 {
		return errors.New("启用/禁用成员: 1表示启用成员，0表示禁用成员")
	}
	if nameLen := utf8.RuneCountInString(a.Name); nameLen < 1 || nameLen > 64 {
		return errors.New("成员名称长度为1~64个utf8字符")
	}
	deptLen := len(a.Department)
	if deptLen > 20 {
		return errors.New("成员所属部门id列表要求不超过20个")
	}
	if deptLen != len(a.IsLeaderInDept) {
		return errors.New("成员上级字段个数必须和department一致")
	}
	for _, isLeader := range a.IsLeaderInDept {
		if isLeader != 0 && isLeader != 1 {
			return errors.New("成员在所在的部门内是否为上级:1表示为上级,0表示非上级")
		}
	}
	if deptLen != len(a.Order) {
		return errors.New("成员排序值字段个数必须和department一致")
	}
	for _, order := range a.Order {
		if order < 0 || order > math.MaxUint32 {
			return errors.New("部门内的排序值有效范围为[0,2^32)")
		}
	}
	if a.Gender != "" && a.Gender != "1" && a.Gender != "2" {
		return errors.New("成员性别必须是1:男,2:女")
	}
	if a.Position != "" && len(a.Position) > 128 {
		return errors.New("成员职务信息长度为0~128个字符")
	}
	if a.Email != "" {
		if len(a.Email) > 64 {
			return errors.New("成员的邮箱地址长度最大为64")
		}
		// 找到@符号的位置
		index := strings.IndexByte(a.Email, '@')
		if index < 1 || index == len(a.Email)-1 {
			return errors.New("成员的email必须是有效的邮箱地址")
		}
	}
	if a.ExternalPosition != "" {
		if utf8.RuneCountInString(a.ExternalPosition) > 12 {
			return errors.New("成员的外部职务必须是最大12个中文字符")
		}
		for _, r := range a.ExternalPosition {
			if !unicode.Is(unicode.Han, r) {
				return errors.New("成员的外部职务必须是最大12个中文字符")
			}
		}
	}
	if a.ExtAttr != nil {
		err := a.ExtAttr.Validate()
		if err != nil {
			return err
		}
	}
	if a.ExternalProfile != nil {
		err := a.ExternalProfile.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// ExternalProfile -
type ExternalProfile struct {
	ExternalCoprName string    `json:"external_corp_name,omitempty" xml:",omitempty"`
	ExternalAttr     []ExtAttr `json:"external_attr,omitempty" xml:",omitempty"`
}

// Validate 验证扩展属性
func (a ExternalProfile) Validate() (err error) {
	for _, attr := range a.ExternalAttr {
		if err = attr.Validate(); err != nil {
			return
		}
	}
	return nil
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
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
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
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
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

// NewCreateUserURL 新建创建用户的URL
func NewCreateUserURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultCreateUserURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

//NewUpdateUserURL 新建更新用户的URL
func NewUpdateUserURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultUpdateUserURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// NewDeleteUserURL 新建删除用户的URL
func NewDeleteUserURL(url, accessToken, userid string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultDeleteUserURL
	}
	return fmt.Sprintf("%s?access_token=%s&userid=%s", url, accessToken, userid)
}

// NewBatchDeleteUserURL 新建批量删除用户的URL
func NewBatchDeleteUserURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultBatchDeleteUserURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// postUser 提交修改用户请求
// 未测试
func postUser(url, accessToken, action string, data interface{}) error {
	if accessToken == "" {
		return errcode.ErrInvalidAccessToken
	}
	switch action {
	case "create", "update":
		_, ok := data.(*User)
		if !ok {
			return errors.New("创建及更新用户操作的数据类型必须是User")
		}

		if action == "create" {
			url = NewCreateUserURL(url, accessToken)
		} else {
			url = NewUpdateUserURL(url, accessToken)
		}
	case "batchdelete":
		_, ok := data.([]string)
		if !ok {
			return errors.New("批量删除用户操作的数据类型必须是用户id数组")
		}

		url = NewBatchDeleteUserURL(url, accessToken)
	default:
		return errors.New("不支持的操作")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(b)
	resp, err := httpClient.Post(url, mimeApplicationJSONCharsetUTF8, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res Response
	if err = json.Unmarshal(b, &res); err != nil {
		return err
	}
	return res.Validate()
}

// CreateUser 创建用户
// 未测试
func CreateUser(url, accessToken string, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return postUser(url, accessToken, "create", user)
}

// UpdateUser 创建用户
// 未测试
func UpdateUser(url, accessToken string, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return postUser(url, accessToken, "update", user)
}

// DeleteUser 删除用户
// 未测试
func DeleteUser(url, accessToken, userid string) error {
	userid = strings.Replace(userid, " ", "", -1)
	if userid == "" {
		return errors.New("成员UserID为空")
	}
	if accessToken == "" {
		return errcode.ErrInvalidAccessToken
	}
	url = NewDeleteUserURL(url, accessToken, userid)
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res Response
	if err = json.Unmarshal(b, &res); err != nil {
		return err
	}
	return res.Validate()
}

// BatchDeleteUser 批量删除用户
// 未测试
func BatchDeleteUser(url, accessToken string, userids []string) error {
	userids = RemoveDuplicateString(userids)
	if len(userids) < 1 {
		return nil
	} else if len(userids) > maxBatchDeleteUserCount {
		return fmt.Errorf("批量删除用户数量上限为%d", maxBatchDeleteUserCount)
	}
	for _, userid := range userids {
		if userid == "" {
			return errors.New("用户id存在空字符串")
		}
	}
	return postUser(url, accessToken, "batchdelete", userids)
}

// SwitchOpenIDAndUserIDResponse userid和openid转换的响应
type SwitchOpenIDAndUserIDResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	OpenID  string `json:"openid,omitempty"`
	UserID  string `json:"userid,omitempty"`
}

// Validate 验证响应
func (res *SwitchOpenIDAndUserIDResponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// NewConverUserIDToOpenIDURL 新建userid转openid的URL
// 未测试
func NewConverUserIDToOpenIDURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultUserIDToOpenIDURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// NewConverOpenIDToUserIDURL 新建openid转userid的URL
// 未测试
func NewConverOpenIDToUserIDURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultOpenIDToUserIDURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// switchOpenIDAndUserID 交换openid和userid
// 未测试
func switchOpenIDAndUserID(url, accessToken, id string, typ int) (s string, err error) {
	if accessToken == "" {
		return "", errcode.ErrInvalidAccessToken
	}
	id = strings.Replace(id, " ", "", -1)

	var buf bytes.Buffer
	switch typ {
	case 1: // userid to openid
		url = NewConverUserIDToOpenIDURL(url, accessToken)
		buf.WriteString(`{"userid":"` + id + `"}`)
	case 2: // openid to userid
		url = NewConverOpenIDToUserIDURL(url, accessToken)
		buf.WriteString(`{"openid":"` + id + `"}`)
	default:
		return "", errors.New("无效的操作类型")
	}

	resp, err := httpClient.Post(url, mimeApplicationJSONCharsetUTF8, &buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res SwitchOpenIDAndUserIDResponse
	if err = json.Unmarshal(b, &res); err != nil {
		return "", err
	}
	if err = res.Validate(); err != nil {
		return "", err
	}
	if typ == 1 {
		s = res.OpenID
	} else {
		s = res.UserID
	}
	return s, nil
}

// ConverUserIDToOpenID userid转openid
// 未测试
func ConverUserIDToOpenID(url, accessToken, userid string) (string, error) {
	if userid == "" {
		return "", errors.New("userid为空")
	}
	return switchOpenIDAndUserID(url, accessToken, userid, 1)
}

// ConverOpenIDToUserID openid转userid
// 未测试
func ConverOpenIDToUserID(url, accessToken, openid string) (string, error) {
	if openid == "" {
		return "", errors.New("openid为空")
	}
	return switchOpenIDAndUserID(url, accessToken, openid, 2)
}
