package payment

// OrderCreateRes 订单创建结果
type OrderCreateRes struct {
	Payment bool   `json:"payment"`           // 是否需要支付
	ID      string `json:"id,omitempty"`      // 订单号
	QRCode  string `json:"qr_code,omitempty"` // 支付二维码指向的地址
}
