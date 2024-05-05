package aria2

// AddURLService 添加URL离线下载服务
type BatchAddURLService struct {
	URLs []string `json:"url" binding:"required"`
	Dst  string   `json:"dst" binding:"required,min=1"`
}

// AddURLService 添加URL离线下载服务
type AddURLService struct {
	URL           string `json:"url" binding:"required"`
	Dst           string `json:"dst" binding:"required,min=1"`
	PreferredNode uint   `json:"preferred_node"`
}
