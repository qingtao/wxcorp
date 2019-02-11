package agent

import (
	"time"

	"github.com/qingtao/wxcorp/corp"
	"github.com/qingtao/wxcorp/corp/errcode"
)

// SendMsg 应用发送消息
func (a *Agent) SendMsg(msg *corp.Msg) error {
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return err
	}
	for i := 0; i < retryTimes; i++ {
		err = corp.SendMsg("", accessToken, msg)
		if err == nil {
			break
		}

		// 只检查令牌错误
		if err != errcode.ErrInvalidAccessToken {
			return err
		}
		accessToken, err = a.RefreshAccessToken()
		if err != nil {
			return err
		}
		time.Sleep(retryInterval)
	}
	return nil
}
