package serializer

type CaptchaReq struct {
	CaptchaCode string `json:"captchaCode"`
	Ticket      string `json:"ticket,omitempty"`
	Randstr     string `json:"randstr,omitempty"`
}
