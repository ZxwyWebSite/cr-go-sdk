package cr

import "net/http"

type UserObj struct {
	Mail string // 邮箱(用户名)
	Pass string // 密码
	// Sess string // 会话
	// Exps int64 // 有效期

	Cookie *http.Cookie
}

func NewUser(mail, pass string) *UserObj {
	return &UserObj{
		Mail: mail,
		Pass: pass,
	}
}
