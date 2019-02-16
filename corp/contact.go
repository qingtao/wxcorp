package corp

// ContactEvent 通讯录变更事件
type ContactEvent struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	AgentID      string
	Event        string
	ChangeType   string `xml:",omitempty"`

	// MsgID 可用于普通消息排重

	// ID 部门id
	ID int `json:"id" xml:"Id,omitempty"`
	// Name 部门名称
	Name string `json:"name"`
	// ParentID 父亲部门id,根部门为1
	ParentID int `json:"parentid" xml:"ParentId"`
	// Order 在父部门中的次序值,order值小的排序靠前
	Order int `json:"order"`

	// UserID 	成员UserID,对应管理端的帐号,企业内部必须唯一，1-64个字节
	UserID    string `json:"userid"`
	NewUserID string `json:"-" xml:",omitempty"`
	// Name 成员名称,1-64个utf8字符
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
	// Telephone 电话号码,32字节以内，可以包含数字和"-"
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

	TagID         int    `xml:"TagId,omitempty"`
	AddUserItems  string `xml:",omitempty"`
	DelUserItems  string `xml:",omitempty"`
	AddPartyItems string `xml:",omitempty"`
	DelPartyItems string `xml:",omitempty"`

	JobID   string `xml:"JobId,omitempty"`
	JobType string `xml:",omitempty"`
	ErrCode int    `xml:",omitempty"`
	ErrMsg  string `xml:",omitempty"`
}

// TagChangeEvent 标签变更事件
type TagChangeEvent struct {
	TagID         int      `xml:"TagId,omitempty"`
	AddUserItems  []string `xml:",omitempty"`
	DelUserItems  []string `xml:",omitempty"`
	AddPartyItems []int    `xml:",omitempty"`
	DelPartyItems []int    `xml:",omitempty"`
}

// BatchJobEvent 异步任务事件
type BatchJobEvent struct {
	JobID   string `xml:"JobId,omitempty"`
	JobType string `xml:",omitempty"`
	ErrCode int    `xml:",omitempty"`
	ErrMsg  string `xml:",omitempty"`
}
