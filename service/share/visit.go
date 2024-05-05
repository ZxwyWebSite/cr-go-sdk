package share

// ShareUserGetService 获取用户的分享服务
type ShareUserGetService struct {
	Type string `form:"type" binding:"required,eq=hot|eq=default"`
	Page uint   `form:"page" binding:"required,min=1"`
}

// ShareGetService 获取分享服务
type ShareGetService struct {
	Password string `form:"password" binding:"max=255"`
}

// Service 对分享进行操作的服务，
// path 为可选文件完整路径，在目录分享下有效
type Service struct {
	Path string `form:"path" uri:"path" binding:"max=65535"`
}

// ArchiveService 分享归档下载服务
type ArchiveService struct {
	Path  string   `json:"path" binding:"required,max=65535"`
	Items []string `json:"items"`
	Dirs  []string `json:"dirs"`
}

// ShareListService 列出分享
type ShareListService struct {
	Page     uint   `form:"page" binding:"required,min=1"`
	OrderBy  string `form:"order_by" binding:"required,eq=created_at|eq=downloads|eq=views"`
	Order    string `form:"order" binding:"required,eq=DESC|eq=ASC"`
	Keywords string `form:"keywords"`
}

// ShareReportService 举报分享
type ShareReportService struct {
	Reason int    `json:"reason" binding:"gte=0,lte=4"`
	Des    string `json:"des"`
}
