package user

// UserRegisterService 管理用户注册的服务
type UserRegisterService struct {
	//TODO 细致调整验证规则
	UserName string `form:"userName" json:"userName" binding:"required,email"`
	Password string `form:"Password" json:"Password" binding:"required,min=4,max=64"`
}
