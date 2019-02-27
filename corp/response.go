package corp

import (
	"github.com/qingtao/wxcorp/corp/errcode"
)

// Response 通用响应
type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// Validate 检查响应信息
func (res *Response) Validate() error {
	if res == nil {
		return ErrIsNil
	}
	return errcode.Error(res.ErrCode)
}
