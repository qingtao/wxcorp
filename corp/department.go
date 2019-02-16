package corp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/qingtao/wxcorp/corp/errcode"
)

const defaultDepartmentListURL = "https://qyapi.weixin.qq.com/cgi-bin/department/list"

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
