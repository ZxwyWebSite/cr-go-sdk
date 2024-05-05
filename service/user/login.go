package user

// UserLoginService 管理用户登录的服务
type UserLoginService struct {
	//TODO 细致调整验证规则
	UserName string `form:"userName" json:"userName" binding:"required,email"`
	Password string `form:"Password" json:"Password" binding:"required,min=4,max=64"`
}

// UserResetEmailService 发送密码重设邮件服务
type UserResetEmailService struct {
	UserName string `form:"userName" json:"userName" binding:"required,email"`
}

// UserResetService 密码重设服务
type UserResetService struct {
	Password string `form:"Password" json:"Password" binding:"required,min=4,max=64"`
	ID       string `json:"id" binding:"required"`
	Secret   string `json:"secret" binding:"required"`
}

// 注册|登录|重置密码 账号信息
type LoginInfo struct {
	UserName    string `form:"userName" json:"userName" binding:"required,email"`
	Password    string `form:"Password" json:"Password" binding:"required,min=4,max=64"`
	CaptchaCode string `json:"captchaCode"`
	Ticket      string `json:"ticket,omitempty"`
	Randstr     string `json:"randstr,omitempty"`
	// UserLoginService
	// serializer.CaptchaReq
}
