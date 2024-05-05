package aria2

// SelectFileService 选择要下载的文件服务
type SelectFileService struct {
	Indexes []int `json:"indexes" binding:"required"`
}

// DownloadTaskService 下载任务管理服务
type DownloadTaskService struct {
	GID string `uri:"gid" binding:"required"`
}

// DownloadListService 下载列表服务
type DownloadListService struct {
	Page uint `form:"page"`
}
