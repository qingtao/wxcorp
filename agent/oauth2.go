package agent

import (
	"time"

	"github.com/qingtao/wxcorp/corp"
	"github.com/qingtao/wxcorp/corp/errcode"
)

// GetUserInfo 通过微信授权返回用户信息
func (a *Agent) GetUserInfo(code string) (userid, deviceid string, err error) {
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return "", "", err
	}
	for i := 0; i < retryTimes; i++ {
		res, err := corp.GetUserInfoWithCode("", accessToken, code)
		// 无错误返回
		if err == nil {
			userid, deviceid = res.UserID, res.DeviceID
			break
		}
		// 如果是令牌错误,主动刷新令牌
		if err == errcode.ErrInvalidAccessToken {
			// 刷新令牌错误返回
			accessToken, err = a.RefreshAccessToken()
			if err != nil {
				return "", "", err
			}
		}
		// 等待指定时间后重试
		time.Sleep(retryInterval)
	}
	return
}
