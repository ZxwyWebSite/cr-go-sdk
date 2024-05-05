package explorer

// SingleFileService 对单文件进行操作的五福，path为文件完整路径
type SingleFileService struct {
	Path string `uri:"path" json:"path" binding:"required,min=1,max=65535"`
}

// FileIDService 通过文件ID对文件进行操作的服务
type FileIDService struct {
}

// FileAnonymousGetService 匿名（外链）获取文件服务
type FileAnonymousGetService struct {
	ID   uint   `uri:"id" binding:"required,min=1"`
	Name string `uri:"name" binding:"required"`
}

// DownloadService 文件下載服务
type DownloadService struct {
	ID string `uri:"id" binding:"required"`
}

// ArchiveService 文件流式打包下載服务
type ArchiveService struct {
	ID string `uri:"sessionID" binding:"required"`
}
