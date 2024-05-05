package explorer

// DirectoryService 创建新目录服务
type DirectoryService struct {
	Path string `uri:"path" json:"path" binding:"required,min=1,max=65535"`
}
