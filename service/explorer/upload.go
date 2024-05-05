package explorer

// CreateUploadSessionService 获取上传凭证服务
type CreateUploadSessionService struct {
	Path         string `json:"path" binding:"required"`
	Size         uint64 `json:"size" binding:"min=0"`
	Name         string `json:"name" binding:"required"`
	PolicyID     string `json:"policy_id" binding:"required"`
	LastModified int64  `json:"last_modified"`
	MimeType     string `json:"mime_type"`
}
