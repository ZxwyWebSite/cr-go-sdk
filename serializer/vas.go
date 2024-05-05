package serializer

import (
	"time"
)

type Quota struct {
	Base  uint64         `json:"base"`
	Pack  uint64         `json:"pack"`
	Used  uint64         `json:"used"`
	Total uint64         `json:"total"`
	Packs []storagePacks `json:"packs"`
}

type storagePacks struct {
	Name           string    `json:"name"`
	Size           uint64    `json:"size"`
	ActivateDate   time.Time `json:"activate_date"`
	Expiration     int       `json:"expiration"`
	ExpirationDate time.Time `json:"expiration_date"`
}

// MountedFolders 已挂载的目录
type MountedFolders struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PolicyName string `json:"policy_name"`
}

type PolicyOptions struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type NodeOptions struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}

// PackProduct 容量包商品
type PackProduct struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Size  uint64 `json:"size"`
	Time  int64  `json:"time"`
	Price int    `json:"price"`
	Score int    `json:"score"`
}

// GroupProducts 用户组商品
type GroupProducts struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	GroupID   uint     `json:"group_id"`
	Time      int64    `json:"time"`
	Price     int      `json:"price"`
	Score     int      `json:"score"`
	Des       []string `json:"des"`
	Highlight bool     `json:"highlight"`
}

// 增值服务商品响应
type ProductData struct {
	Packs      []PackProduct   `json:"packs"`
	Groups     []GroupProducts `json:"groups"`
	Alipay     bool            `json:"alipay"`
	Wechat     bool            `json:"wechat"`
	Payjs      bool            `json:"payjs"`
	Custom     bool            `json:"custom"`
	CustomName string          `json:"custom_name"`
	ScorePrice int             `json:"score_price"`
}
