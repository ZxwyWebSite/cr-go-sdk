package serializer

import (
	"time"
)

// SiteConfig 站点全局设置序列
type SiteConfig struct {
	SiteName             string   `json:"title"`
	LoginCaptcha         bool     `json:"loginCaptcha"`
	RegCaptcha           bool     `json:"regCaptcha"`
	ForgetCaptcha        bool     `json:"forgetCaptcha"`
	EmailActive          bool     `json:"emailActive"`
	QQLogin              bool     `json:"QQLogin"`
	Themes               string   `json:"themes"`
	DefaultTheme         string   `json:"defaultTheme"`
	ScoreEnabled         bool     `json:"score_enabled"`
	ShareScoreRate       string   `json:"share_score_rate"`
	HomepageViewMethod   string   `json:"home_view_method"`
	ShareViewMethod      string   `json:"share_view_method"`
	Authn                bool     `json:"authn"`
	User                 User     `json:"user"`
	ReCaptchaKey         string   `json:"captcha_ReCaptchaKey"`
	SiteNotice           string   `json:"site_notice"`
	CaptchaType          string   `json:"captcha_type"`
	TCaptchaCaptchaAppId string   `json:"tcaptcha_captcha_app_id"`
	RegisterEnabled      bool     `json:"registerEnabled"`
	ReportEnabled        bool     `json:"report_enabled"`
	AppPromotion         bool     `json:"app_promotion"`
	WopiExts             []string `json:"wopi_exts"`
	AppFeedbackLink      string   `json:"app_feedback"`
	AppForumLink         string   `json:"app_forum"`
}

type Task struct {
	Status     int       `json:"status"`
	Type       int       `json:"type"`
	CreateDate time.Time `json:"create_date"`
	Progress   int       `json:"progress"`
	Error      string    `json:"error"`
}

// 任务列表响应
type TaskList struct {
	Total int    `json:"total"`
	Tasks []Task `json:"tasks"`
}

// VolResponse VOL query response
type VolResponse struct {
	Signature string `json:"signature"`
	Content   string `json:"content"`
}
