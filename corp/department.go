package corp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/qingtao/wxcorp/corp/errcode"
)

const (
	defaultDepartmentListURL   = "https://qyapi.weixin.qq.com/cgi-bin/department/list"
	defaultDepartmentCreateURL = "https://qyapi.weixin.qq.com/cgi-bin/department/create"
	defaultDepartmentUpdateURL = "https://qyapi.weixin.qq.com/cgi-bin/department/update"
	defaultDepartmentDeleteURL = "https://qyapi.weixin.qq.com/cgi-bin/department/delete"
)

// Department 部门结构
type Department struct {
	// ID 部门id
	ID int `json:"id"`
	// Name 部门名称
	Name string `json:"name"`
	// ParentID 父亲部门id,根部门为1
	ParentID int `json:"parentid"`
	// Order 在父部门中的次序值,order值大的排序靠前
	Order int `json:"order"`
}

var (
	deptInvalidName   = `[\:?"<>'/]`
	reDeptInvalidName = regexp.MustCompile(deptInvalidName)
)

// Validate 验证部门的参数是否符合企业微信的规则
func (a *Department) Validate(action string) error {
	// 创建时name必须设置
	if action == "create" && a.Name == "" {
		return errors.New("部门名称长度限制为1-32个字符")
	}
	// 设置name时校验名称
	if a.Name != "" {
		a.Name = strings.Replace(a.Name, " ", "", -1)
		if len(a.Name) > 32 {
			return errors.New("部门名称长度限制为1-32个字符")
		}
		if reDeptInvalidName.MatchString(a.Name) {
			return errors.Errorf(`部门名称不能包含%s`, deptInvalidName)
		}
	}
	if a.ParentID == a.ID {
		return errors.New("部门ID和父部门ID相同")
	}
	if a.ParentID < 0 || a.ParentID > math.MaxUint32 {
		return errors.New("父部门ID有效范围[1,2^32)")
	}
	if a.Order < 0 || a.Order > math.MaxUint32 {
		return errors.New("部门的排序值有效范围是[0,2^32)")
	}
	if a.ID < 1 || a.ID > math.MaxUint32 {
		return errors.New("部门ID指定时必须大于1,不指定时自动生成,有效范围[1,2^32)")
	}
	return nil
}

// ChangeDepartmentResponse 创建部门的响应结构
type ChangeDepartmentResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	ID      int    `json:"id"`
}

// Validate 验证创建部门的响应数据
func (res *ChangeDepartmentResponse) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}

// DepartmentResponse 请求部门列表的响应
type DepartmentResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// Department 部门列表数据,以部门的order字段从大到小排列
	Department []Department `json:"department"`
}

// Validate 检查响应
func (dept *DepartmentResponse) Validate() error {
	if dept == nil {
		return ErrIsNil
	}
	return errcode.Error(dept.ErrCode)
}

// NewGetDepartmentListURL 新建获取部门列表的URL
func NewGetDepartmentListURL(url, accessToken string, id int) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultDepartmentListURL
	}
	return fmt.Sprintf("%s?access_token=%s&id=%d", url, accessToken, id)
}

// GetDepartment 获取部门信息
// 未测试
func GetDepartment(url, accessToken string, id int) (dept *DepartmentResponse, err error) {
	if accessToken == "" {
		return nil, errcode.ErrInvalidAccessToken
	}
	url = NewGetDepartmentListURL(url, accessToken, id)
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
	if err = json.Unmarshal(b, &dept); err != nil {
		return nil, err
	}
	err = dept.Validate()
	return
}

// NewCreateDepartmentURL 新建创建部门请求的URL
func NewCreateDepartmentURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultDepartmentCreateURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// NewUpdateDepartmentURL 新建更新部门请求的URL
func NewUpdateDepartmentURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultDepartmentUpdateURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// postDepartment 提交部门, action指定为"create"时按照创建部门检查名称是否为空
// 未测试
func postDepartment(url, accessToken, action string, dept *Department) (err error) {
	if accessToken == "" {
		return errcode.ErrInvalidAccessToken
	}
	switch action {
	case "create":
		err = dept.Validate(action)
		if err != nil {
			return
		}
		url = NewCreateDepartmentURL(url, accessToken)
	case "update":
		err = dept.Validate(action)
		if err != nil {
			return
		}
		url = NewUpdateDepartmentURL(url, accessToken)
	default:
		return errors.New("不支持的操作")
	}
	b, err := json.Marshal(dept)
	if err != nil {
		return
	}
	buf := bytes.NewReader(b)
	resp, err := httpClient.Post(url, mimeApplicationJSONCharsetUTF8, buf)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var res ChangeDepartmentResponse
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}
	err = res.Validate()
	return
}

// CreateDepartment 创建部门
// 未测试
func CreateDepartment(url, accessToken string, dept *Department) error {
	return postDepartment(url, accessToken, "create", dept)
}

// UpdateDepartment 更新部门
// 未测试
func UpdateDepartment(url, accessToken string, dept *Department) error {
	return postDepartment(url, accessToken, "update", dept)
}

// NewDeleteDepartmentURL 新建删除部门的URL
func NewDeleteDepartmentURL(url, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	if url == "" {
		url = defaultDepartmentDeleteURL
	}
	return fmt.Sprintf("%s?access_token=%s", url, accessToken)
}

// DeleteDepartment 删除部门,不能删除根部门；不能删除含有子部门、成员的部门, 所以要先确定id是否可以删除
// 未测试
func DeleteDepartment(url, accessToken string, id int) error {
	if accessToken == "" {
		return errcode.ErrInvalidAccessToken
	}
	url = NewDeleteDepartmentURL(url, accessToken)
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res ChangeDepartmentResponse
	if err = json.Unmarshal(b, &res); err != nil {
		return err
	}
	return res.Validate()
}
