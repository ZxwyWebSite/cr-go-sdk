package explorer

// FilterTagCreateService 文件分类标签创建服务
type FilterTagCreateService struct {
	Expression string `json:"expression" binding:"required,min=1,max=65535"`
	Icon       string `json:"icon" binding:"required,min=1,max=255"`
	Name       string `json:"name" binding:"required,min=1,max=255"`
	Color      string `json:"color" binding:"hexcolor|rgb|rgba|hsl"`
}

// LinkTagCreateService 目录快捷方式标签创建服务
type LinkTagCreateService struct {
	Path string `json:"path" binding:"required,min=1,max=65535"`
	Name string `json:"name" binding:"required,min=1,max=255"`
}
