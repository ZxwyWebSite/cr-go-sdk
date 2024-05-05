package user

import (
	"time"

	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
)

// SettingListService 通用设置列表服务
type SettingListService struct {
	Page int `form:"page" binding:"required,min=1"`
}

// AvatarService 头像服务
type AvatarService struct {
	Size string `uri:"size" binding:"required,eq=l|eq=m|eq=s"`
}

// SettingUpdateService 设定更改服务
type SettingUpdateService struct {
	Option string `uri:"option" binding:"required,eq=nick|eq=theme|eq=homepage|eq=vip|eq=qq|eq=policy|eq=password|eq=2fa|eq=authn"`
}

// ChangerNick 昵称更改服务
type ChangerNick struct {
	Nick string `json:"nick" binding:"required,min=1,max=255"`
}

// VIPUnsubscribe 用户组解约服务
type VIPUnsubscribe struct {
}

// QQBind QQ互联服务
type QQBind struct {
}

// PolicyChange 更改存储策略
type PolicyChange struct {
	ID string `json:"id" binding:"required"`
}

// HomePage 更改个人主页开关
type HomePage struct {
	Enabled bool `json:"status"`
}

// PasswordChange 更改密码
type PasswordChange struct {
	Old string `json:"old" binding:"required,min=4,max=64"`
	New string `json:"new" binding:"required,min=4,max=64"`
}

// Enable2FA 开启二步验证
type Enable2FA struct {
	Code string `json:"code" binding:"required"`
}

// DeleteWebAuthn 删除WebAuthn凭证
type DeleteWebAuthn struct {
	ID string `json:"id" binding:"required"`
}

// ThemeChose 主题选择
type ThemeChose struct {
	Theme string `json:"theme" binding:"required,hexcolor|rgb|rgba|hsl"`
}

// 用户设定
type Settings struct {
	Uid          uint                             `json:"uid"`
	QQ           bool                             `json:"qq"`
	HomePage     bool                             `json:"homepage"`
	TwoFactor    bool                             `json:"two_factor"`
	PreferTheme  string                           `json:"prefer_theme"`
	Themes       string                           `json:"themes"`
	GroupExpires *time.Time                       `json:"group_expires"`
	Authn        []serializer.WebAuthnCredentials `json:"authn"`
}

// 更改用户设定
type UpdOption struct {
	ChangerNick
	VIPUnsubscribe
	QQBind
	PolicyChange
	HomePage
	PasswordChange
	Enable2FA
	DeleteWebAuthn
	ThemeChose
}
