package share

// ShareCreateService 创建新分享服务
type ShareCreateService struct {
	SourceID        string `json:"id" binding:"required"`
	IsDir           bool   `json:"is_dir"`
	Password        string `json:"password" binding:"max=255"`
	RemainDownloads int    `json:"downloads"`
	Expire          int    `json:"expire"`
	Score           int    `json:"score" binding:"gte=0"`
	Preview         bool   `json:"preview"`
}

// ShareUpdateService 分享更新服务
type ShareUpdateService struct {
	Prop  string `json:"prop" binding:"required,eq=password|eq=preview_enabled"`
	Value string `json:"value" binding:"max=255"`
}
