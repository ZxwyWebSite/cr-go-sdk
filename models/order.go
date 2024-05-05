package models

const (
	// PackOrderType 容量包订单
	PackOrderType = iota
	// GroupOrderType 用户组订单
	GroupOrderType
	// ScoreOrderType 积分充值订单
	ScoreOrderType
)

const (
	// OrderUnpaid 未支付
	OrderUnpaid = iota
	// OrderPaid 已支付
	OrderPaid
	// OrderCanceled 已取消
	OrderCanceled
)
