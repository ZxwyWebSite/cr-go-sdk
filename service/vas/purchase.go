package vas

// CreateOrderService 创建订单服务
type CreateOrderService struct {
	Action string `json:"action" binding:"required,eq=group|eq=pack|eq=score"`
	Method string `json:"method" binding:"required,eq=alipay|eq=score|eq=payjs|eq=wechat|eq=custom"`
	ID     int64  `json:"id" binding:"required"`
	Num    int    `json:"num" binding:"required,min=1"`
}

// RedeemService 兑换服务
type RedeemService struct {
	Code string `uri:"code" binding:"required,max=64"`
}

// OrderService 订单查询
type OrderService struct {
	ID string `uri:"id" binding:"required"`
}

// 兑换码信息
type RedeemData struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Num  int    `json:"num"`
	Time int64  `json:"time"`
}
