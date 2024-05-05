package explorer

import "github.com/ZxwyWebSite/cr-go-sdk/serializer"

// ItemSearchService 文件搜索服务
type ItemSearchService struct {
	Type     string `uri:"type" binding:"required"`
	Keywords string `uri:"keywords" binding:"required"`
	Path     string `form:"path"`
}

// 根据关键字搜索文件结果
type SearchResult struct {
	Parent  int                 `json:"parent"`
	Objects []serializer.Object `json:"objects"`
}
