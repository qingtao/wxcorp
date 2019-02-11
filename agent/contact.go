package agent

import (
	"time"

	"github.com/qingtao/wxcorp/corp"
	"github.com/qingtao/wxcorp/corp/errcode"
)

// GetDepartment 获取部门列表
func (a *Agent) GetDepartment(departmentID int) (dept []corp.Department, err error) {
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return nil, err
	}
	for i := 0; i < retryTimes; i++ {
		res, err := corp.GetDepartment("", accessToken, departmentID)
		if err == nil {
			dept = res.Department
			break
		}
		// 如果是未重试且错误是令牌问题，则等待[retryInterval]秒刷新令牌并尝试重试一次
		if err == errcode.ErrInvalidAccessToken {
			accessToken, err = a.RefreshAccessToken()
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(retryInterval)
	}
	return
}

// GetUserListOfDepartment 获取部门成员
func (a *Agent) GetUserListOfDepartment(departmentID, fetchChild int, typ string) (users []corp.User, err error) {
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return nil, err
	}
	for i := 0; i < retryTimes; i++ {
		res, err := corp.GetUserList("", typ, accessToken, departmentID, fetchChild, nil)
		if err == nil {
			users = res.UserList
			break
		}
		if err == errcode.ErrInvalidAccessToken {
			accessToken, err = a.RefreshAccessToken()
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(retryInterval)
	}
	return
}

// GetTagList 获取标签列表
func (a *Agent) GetTagList() (tags []corp.Tag, err error) {
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return nil, err
	}
	for i := 0; i < retryTimes; i++ {
		res, err := corp.GetTagList("", accessToken)
		if err == nil {
			tags = res.TagList
			break
		}
		if err == errcode.ErrInvalidAccessToken {
			accessToken, err = a.RefreshAccessToken()
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(retryInterval)
	}
	return
}

// GetMemberOfTag 获取标签成员
func (a *Agent) GetMemberOfTag(id int) (member *corp.Member, err error) {
	accessToken, err := a.GetAccessToken()
	if err != nil {
		return nil, err
	}
	for i := 0; i < retryTimes; i++ {
		res, err := corp.GetMemberOfTag("", accessToken, id)
		if err == nil {
			member = &corp.Member{
				TagName:   res.TagName,
				UserList:  res.UserList,
				PartyList: res.PartyList,
			}
			break
		}
		if err == errcode.ErrInvalidAccessToken {
			accessToken, err = a.RefreshAccessToken()
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(retryInterval)
	}
	return
}
