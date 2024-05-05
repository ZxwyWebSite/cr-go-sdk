package serializer

import "strconv"

// AppError 应用错误，实现了error接口
type AppError struct {
	Code int
	Msg  string
	// RawError error
}

// NewError 返回新的错误对象
/*func NewError(code int, msg string, err error) AppError {
	return AppError{
		Code:     code,
		Msg:      msg,
		RawError: err,
	}
}

// NewErrorFromResponse 从 serializer.Response 构建错误
func NewErrorFromResponse(resp *Response) AppError {
	return AppError{
		Code:     resp.Code,
		Msg:      resp.Msg,
		RawError: errors.New(resp.Error),
	}
}

// WithError 将应用error携带标准库中的error
func (err *AppError) WithError(raw error) AppError {
	err.RawError = raw
	return *err
}*/

// Error 返回业务代码确定的可读错误信息
func (err AppError) Error() string {
	if err.Msg == `` {
		if e, ok := CodeMap[err.Code]; ok {
			err.Msg = e
		}
	}
	return strconv.Itoa(err.Code) + `: ` + err.Msg
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	// CodeNotFullySuccess 未完全成功
	CodeNotFullySuccess = 203
	// CodeCheckLogin 未登录
	CodeCheckLogin = 401
	// CodeNoPermissionErr 未授权访问
	CodeNoPermissionErr = 403
	// CodeNotFound 资源未找到
	CodeNotFound = 404
	// CodeConflict 资源冲突
	CodeConflict = 409
	// CodeUploadFailed 上传出错
	CodeUploadFailed = 40002
	// CodeCreateFolderFailed 目录创建失败
	CodeCreateFolderFailed = 40003
	// CodeObjectExist 对象已存在
	CodeObjectExist = 40004
	// CodeSignExpired 签名过期
	CodeSignExpired = 40005
	// CodePolicyNotAllowed 当前存储策略不允许
	CodePolicyNotAllowed = 40006
	// CodeGroupNotAllowed 用户组无法进行此操作
	CodeGroupNotAllowed = 40007
	// CodeAdminRequired 非管理用户组
	CodeAdminRequired = 40008
	// CodeMasterNotFound 主机节点未注册
	CodeMasterNotFound = 40009
	// CodeUploadSessionExpired 上传会话已过期
	CodeUploadSessionExpired = 40011
	// CodeInvalidChunkIndex 无效的分片序号
	CodeInvalidChunkIndex = 40012
	// CodeInvalidContentLength 无效的正文长度
	CodeInvalidContentLength = 40013
	// CodePhoneRequired 未绑定手机
	CodePhoneRequired = 40010
	// CodeBatchSourceSize 超出批量获取外链限制
	CodeBatchSourceSize = 40014
	// CodeBatchAria2Size 超出最大 Aria2 任务数量限制
	CodeBatchAria2Size = 40015
	// CodeParentNotExist 父目录不存在
	CodeParentNotExist = 40016
	// CodeUserBaned 用户不活跃
	CodeUserBaned = 40017
	// CodeUserNotActivated 用户不活跃
	CodeUserNotActivated = 40018
	// CodeFeatureNotEnabled 此功能未开启
	CodeFeatureNotEnabled = 40019
	// CodeCredentialInvalid 凭证无效
	CodeCredentialInvalid = 40020
	// CodeUserNotFound 用户不存在
	CodeUserNotFound = 40021
	// Code2FACodeErr 二步验证代码错误
	Code2FACodeErr = 40022
	// CodeLoginSessionNotExist 登录会话不存在
	CodeLoginSessionNotExist = 40023
	// CodeInitializeAuthn 无法初始化 WebAuthn
	CodeInitializeAuthn = 40024
	// CodeWebAuthnCredentialError WebAuthn 凭证无效
	CodeWebAuthnCredentialError = 40025
	// CodeCaptchaError 验证码错误
	CodeCaptchaError = 40026
	// CodeCaptchaRefreshNeeded 验证码需要刷新
	CodeCaptchaRefreshNeeded = 40027
	// CodeFailedSendEmail 邮件发送失败
	CodeFailedSendEmail = 40028
	// CodeInvalidTempLink 临时链接无效
	CodeInvalidTempLink = 40029
	// CodeTempLinkExpired 临时链接过期
	CodeTempLinkExpired = 40030
	// CodeEmailProviderBaned 邮箱后缀被禁用
	CodeEmailProviderBaned = 40031
	// CodeEmailExisted 邮箱已被使用
	CodeEmailExisted = 40032
	// CodeEmailSent 邮箱已重新发送
	CodeEmailSent = 40033
	// CodeUserCannotActivate 用户无法激活
	CodeUserCannotActivate = 40034
	// 存储策略不存在
	CodePolicyNotExist = 40035
	// 无法删除默认存储策略
	CodeDeleteDefaultPolicy = 40036
	// 存储策略下还有文件
	CodePolicyUsedByFiles = 40037
	// 存储策略绑定了用户组
	CodePolicyUsedByGroups = 40038
	// 用户组不存在
	CodeGroupNotFound = 40039
	// 对系统用户组执行非法操作
	CodeInvalidActionOnSystemGroup = 40040
	// 用户组正在被使用
	CodeGroupUsedByUser = 40041
	// 为初始用户更改用户组
	CodeChangeGroupForDefaultUser = 40042
	// 对系统用户执行非法操作
	CodeInvalidActionOnDefaultUser = 40043
	// 文件不存在
	CodeFileNotFound = 40044
	// 列取文件失败
	CodeListFilesError = 40045
	// 对系统节点进行非法操作
	CodeInvalidActionOnSystemNode = 40046
	// 创建文件系统出错
	CodeCreateFSError = 40047
	// 创建任务出错
	CodeCreateTaskError = 40048
	// 文件尺寸太大
	CodeFileTooLarge = 40049
	// 文件类型不允许
	CodeFileTypeNotAllowed = 40050
	// 用户容量不足
	CodeInsufficientCapacity = 40051
	// 对象名非法
	CodeIllegalObjectName = 40052
	// 不支持对根目录执行此操作
	CodeRootProtected = 40053
	// 当前目录下已经有同名文件正在上传中
	CodeConflictUploadOngoing = 40054
	// 文件信息不一致
	CodeMetaMismatch = 40055
	// 不支持该格式的压缩文件
	CodeUnsupportedArchiveType = 40056
	// 可用存储策略发生变化
	CodePolicyChanged = 40057
	// 分享链接无效
	CodeShareLinkNotFound = 40058
	// 不能转存自己的分享
	CodeSaveOwnShare = 40059
	// 从机无法向主机发送回调请求
	CodeSlavePingMaster = 40060
	// Cloudreve 版本不一致
	CodeVersionMismatch = 40061
	// 积分不足
	CodeInsufficientCredit = 40062
	// 用户组冲突
	CodeGroupConflict = 40063
	// 当前已处于此用户组中
	CodeGroupInvalid = 40064
	// 兑换码无效
	CodeInvalidGiftCode = 40065
	// 已绑定了QQ账号
	CodeQQBindConflict = 40066
	// QQ账号已被绑定其他账号
	CodeQQBindOtherAccount = 40067
	// QQ 未绑定对应账号
	CodeQQNotLinked = 40068
	// 密码不正确
	CodeIncorrectPassword = 40069
	// 分享无法预览
	CodeDisabledSharePreview = 40070
	// 签名无效
	CodeInvalidSign = 40071
	// 管理员无法购买用户组
	CodeFulfillAdminGroup = 40072
	// CodeDBError 数据库操作失败
	CodeDBError = 50001
	// CodeEncryptError 加密失败
	CodeEncryptError = 50002
	// CodeIOFailed IO操作失败
	CodeIOFailed = 50004
	// CodeInternalSetting 内部设置参数错误
	CodeInternalSetting = 50005
	// CodeCacheOperation 缓存操作失败
	CodeCacheOperation = 50006
	// CodeCallbackError 回调失败
	CodeCallbackError = 50007
	// 后台设置更新失败
	CodeUpdateSetting = 50008
	// 跨域策略添加失败
	CodeAddCORS = 50009
	// 节点不可用
	CodeNodeOffline = 50010
	// 文件元信息查询失败
	CodeQueryMetaFailed = 50011
	//CodeParamErr 各种奇奇怪怪的参数错误
	CodeParamErr = 40001
	// CodeNotSet 未定错误，后续尝试从error中获取
	CodeNotSet = -1
)

var CodeMap = map[int]string{
	CodeNotFullySuccess:            `未完全成功`,
	CodeCheckLogin:                 `未登录`,
	CodeNoPermissionErr:            `未授权访问`,
	CodeNotFound:                   `资源未找到`,
	CodeConflict:                   `资源冲突`,
	CodeUploadFailed:               `上传出错`,
	CodeCreateFolderFailed:         `目录创建失败`,
	CodeObjectExist:                `对象已存在`,
	CodeSignExpired:                `签名过期`,
	CodePolicyNotAllowed:           `当前存储策略不允许`,
	CodeGroupNotAllowed:            `用户组无法进行此操作`,
	CodeAdminRequired:              `非管理用户组`,
	CodeMasterNotFound:             `主机节点未注册`,
	CodeUploadSessionExpired:       `上传会话已过期`,
	CodeInvalidChunkIndex:          `无效的分片序号`,
	CodeInvalidContentLength:       `无效的正文长度`,
	CodePhoneRequired:              `未绑定手机`,
	CodeBatchSourceSize:            `超出批量获取外链限制`,
	CodeBatchAria2Size:             `超出最大 Aria2 任务数量限制`,
	CodeParentNotExist:             `父目录不存在`,
	CodeUserBaned:                  `用户被封禁`,
	CodeUserNotActivated:           `用户不活跃`,
	CodeFeatureNotEnabled:          `此功能未开启`,
	CodeCredentialInvalid:          `凭证无效`,
	CodeUserNotFound:               `用户不存在`,
	Code2FACodeErr:                 `二步验证代码错误`,
	CodeLoginSessionNotExist:       `登录会话不存在`,
	CodeInitializeAuthn:            `无法初始化 WebAuthn`,
	CodeWebAuthnCredentialError:    `WebAuthn 凭证无效`,
	CodeCaptchaError:               `验证码错误`,
	CodeCaptchaRefreshNeeded:       `验证码需要刷新`,
	CodeFailedSendEmail:            `邮件发送失败`,
	CodeInvalidTempLink:            `临时链接无效`,
	CodeTempLinkExpired:            `临时链接过期`,
	CodeEmailProviderBaned:         `邮箱后缀被禁用`,
	CodeEmailExisted:               `邮箱已被使用`,
	CodeEmailSent:                  `邮箱已重新发送`,
	CodeUserCannotActivate:         `用户无法激活`,
	CodePolicyNotExist:             `存储策略不存在`,
	CodeDeleteDefaultPolicy:        `无法删除默认存储策略`,
	CodePolicyUsedByFiles:          `存储策略下还有文件`,
	CodePolicyUsedByGroups:         `存储策略绑定了用户组`,
	CodeGroupNotFound:              `用户组不存在`,
	CodeInvalidActionOnSystemGroup: `对系统用户组执行非法操作`,
	CodeGroupUsedByUser:            `用户组正在被使用`,
	CodeChangeGroupForDefaultUser:  `为初始用户更改用户组`,
	CodeInvalidActionOnDefaultUser: `对系统用户执行非法操作`,
	CodeFileNotFound:               `文件不存在`,
	CodeListFilesError:             `列取文件失败`,
	CodeInvalidActionOnSystemNode:  `对系统节点进行非法操作`,
	CodeCreateFSError:              `创建文件系统出错`,
	CodeCreateTaskError:            `创建任务出错`,
	CodeFileTooLarge:               `文件尺寸太大`,
	CodeFileTypeNotAllowed:         `文件类型不允许`,
	CodeInsufficientCapacity:       `用户容量不足`,
	CodeIllegalObjectName:          `对象名非法`,
	CodeRootProtected:              `不支持对根目录执行此操作`,
	CodeConflictUploadOngoing:      `当前目录下已经有同名文件正在上传中`,
	CodeMetaMismatch:               `文件信息不一致`,
	CodeUnsupportedArchiveType:     `不支持该格式的压缩文件`,
	CodePolicyChanged:              `可用存储策略发生变化`,
	CodeShareLinkNotFound:          `分享链接无效`,
	CodeSaveOwnShare:               `不能转存自己的分享`,
	CodeSlavePingMaster:            `从机无法向主机发送回调请求`,
	CodeVersionMismatch:            `Cloudreve 版本不一致`,
	CodeInsufficientCredit:         `积分不足`,
	CodeGroupConflict:              `用户组冲突`,
	CodeGroupInvalid:               `当前已处于此用户组中`,
	CodeInvalidGiftCode:            `兑换码无效`,
	CodeQQBindConflict:             `已绑定了QQ账号`,
	CodeQQBindOtherAccount:         `QQ账号已被绑定其他账号`,
	CodeQQNotLinked:                `QQ 未绑定对应账号`,
	CodeIncorrectPassword:          `密码不正确`,
	CodeDisabledSharePreview:       `分享无法预览`,
	CodeInvalidSign:                `签名无效`,
	CodeFulfillAdminGroup:          `管理员无法购买用户组`,
	CodeDBError:                    `数据库操作失败`,
	CodeEncryptError:               `加密失败`,
	CodeIOFailed:                   `IO操作失败`,
	CodeInternalSetting:            `内部设置参数错误`,
	CodeCacheOperation:             `缓存操作失败`,
	CodeCallbackError:              `回调失败`,
	CodeUpdateSetting:              `后台设置更新失败`,
	CodeAddCORS:                    `跨域策略添加失败`,
	CodeNodeOffline:                `节点不可用`,
	CodeQueryMetaFailed:            `文件元信息查询失败`,
	CodeParamErr:                   `各种奇奇怪怪的参数错误`,
	CodeNotSet:                     `未定错误，后续尝试从error中获取`,
}
