package setting

import (
	"time"

	"github.com/ZxwyWebSite/cr-go-sdk/models"
	"github.com/ZxwyWebSite/cr-go-sdk/serializer"
)

// WebDAVAccountService WebDAV 账号管理服务
type WebDAVAccountService struct {
	ID uint `json:"id" uri:"id" binding:"required,min=1"`
}

// WebDAVAccountCreateService WebDAV 账号创建服务
type WebDAVAccountCreateService struct {
	Path string `json:"path" binding:"required,min=1,max=65535"`
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// WebDAVAccountUpdateService WebDAV 修改只读性和是否使用代理服务
type WebDAVAccountUpdateService struct {
	ID       uint  `json:"id" binding:"required,min=1"`
	Readonly *bool `json:"readonly" binding:"required_without=UseProxy"`
	UseProxy *bool `json:"use_proxy" binding:"required_without=Readonly"`
}

// WebDAVAccountUpdateReadonlyService WebDAV 修改只读性服务
type WebDAVAccountUpdateReadonlyService struct {
	ID       uint `json:"id" binding:"required,min=1"`
	Readonly bool `json:"readonly"`
}

// WebDAVMountCreateService WebDAV 挂载创建服务
type WebDAVMountCreateService struct {
	Path   string `json:"path" binding:"required,min=1,max=65535"`
	Policy string `json:"policy" binding:"required,min=1"`
}

// WebDAV账号列表
type WebDAVAccountList struct {
	Accounts []models.Webdav             `json:"accounts"`
	Folders  []serializer.MountedFolders `json:"folders"`
}

// WebDAV账户
type WebDAVAccountData struct {
	ID        uint      `json:"id"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
