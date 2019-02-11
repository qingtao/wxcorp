// Package ratelimit 描述了微信企业号接口相关的访问频率限制
//	以上所有频率，按天拦截则被屏蔽一天（自然天），按月拦截则屏蔽一个月（30天，非自然月），按分钟拦截则被屏蔽60秒，按小时拦截则被屏蔽60分钟
package ratelimit

const (
	// 基础频率

	// TimesCgiAPIOneCorpOneMinute 每企业调用单个cgi/api不可超过2000次/分，30000次/小时
	TimesCgiAPIOneCorpOneMinute = 2000
	// TimesCgiAPIOneCorpOneHour -
	TimesCgiAPIOneCorpOneHour = 30000
	// TimesCgiAPIOneIPOneMinute 企业每ip调用单个cgi/api不可超过20000次/分，600000次/小时
	TimesCgiAPIOneIPOneMinute = 20000
	// TimesCgiAPIOneIPOneHour -
	TimesCgiAPIOneIPOneHour = 600000

	// 发送应用频率, 发消息频率不计入基础频率

	// BaseTimesSendMsgOneCorp 每企业不可超过帐号上限数*30人次/天（注：若调用api一次发给1000人，算1000人次；若企业帐号上限是500人，则每天可发送15000人次的消息
	BaseTimesSendMsgOneCorp = 30
	// BaseTimesSendMsgOneAppUserMinute 每应用对同一个成员不可超过30条/分，超过部分会被丢弃不下发
	BaseTimesSendMsgOneAppUserMinute = 30

	// BaseTimesCreateAccountMonth 每企业创建帐号数不可超过帐号上限数*3/月
	BaseTimesCreateAccountMonth = 3
	// BaseTimesCreateAppMonth 每企业最大应用数限制为：人员上限/100个（最低30个，最高300个）；创建应用次数不可超过最大应用数*3/月
	BaseTimesCreateAppMonth = 3
	// BaseTimesCreateAppUsersDivisor -
	BaseTimesCreateAppUsersDivisor = 100
	// TimesCreateAppMini -
	TimesCreateAppMini = 30
	// TimesCreateAppMax -
	TimesCreateAppMax = 300
	// TimesSettingDomain 每企业设置的可信域名数不可超过20/月
	TimesSettingDomain = 20
	// TimesGetCheckInData 每企业调用打卡数据getcheckindata不可超过600次/分
	TimesGetCheckInData = 600
	// TimesGetAppRovalData 每企业调用审批数据频率getapprovaldata不可超过600次/分
	TimesGetAppRovalData = 600
	// BaseTimesVisitNewsExternal 图文访问频率 每篇文章被企业外部人员阅读的次数不能超过企业成员上限*3（假设企业成员上限是200，则一篇文章被企业外部人员阅读次数不能超过600次，超过后外部人员无法打开，企业内成员不受影响。）
	BaseTimesVisitNewsExternal = 3
	// TimesJsAPITicketOneCorpHour 获取jsapi_ticket频率 一小时内，一个企业最多可获取400次，且单个应用不能超过100次
	TimesJsAPITicketOneCorpHour = 400
	// TimesJsAPITicketOneCorpAppHour -
	TimesJsAPITicketOneCorpAppHour = 100
	// TimesAppCreateChatGroupOneCorpDay 每企业所有应用创建群数累积不可超过1000/天
	TimesAppCreateChatGroupOneCorpDay = 1000
	// TimesAppCreateChatGroupOneCorpHour 每企业所有应用变更群次数累积不可超过100/小时
	TimesAppCreateChatGroupOneCorpHour = 100
	// TimesSendMsgToChatGroupMax 每企业所有应用发送群消息不可超过20000人次/分，不可超过200000人次/小时（若群有100人，每发一次消息算100人次）
	TimesSendMsgToChatGroupMax = 20000
	// TimesRecieveMsgFromChatGroupMax 每个成员在群中收到的应用消息不可超过200条/分，1万条/天，超过会被丢弃（接口不会报错）
	TimesRecieveMsgFromChatGroupMax = 200
	// TimesUpdloadImage 每个企业，每天最多可以上传100张永久图片，接口详情参见
	TimesUpdloadImage = 100
)
