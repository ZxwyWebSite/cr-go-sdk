package serializer

import (
	"time"
)

// Share 分享信息序列化
type Share struct {
	Key        string        `json:"key"`
	Locked     bool          `json:"locked"`
	IsDir      bool          `json:"is_dir"`
	Score      int           `json:"score"`
	CreateDate time.Time     `json:"create_date,omitempty"`
	Downloads  int           `json:"downloads"`
	Views      int           `json:"views"`
	Expire     int64         `json:"expire"`
	Preview    bool          `json:"preview"`
	Creator    *shareCreator `json:"creator,omitempty"`
	Source     *shareSource  `json:"source,omitempty"`
}

type shareCreator struct {
	Key       string `json:"key"`
	Nick      string `json:"nick"`
	GroupName string `json:"group_name"`
}

type shareSource struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

// myShareItem 我的分享列表条目
type myShareItem struct {
	Key             string       `json:"key"`
	IsDir           bool         `json:"is_dir"`
	Score           int          `json:"score"`
	Password        string       `json:"password"`
	CreateDate      time.Time    `json:"create_date,omitempty"`
	Downloads       int          `json:"downloads"`
	RemainDownloads int          `json:"remain_downloads"`
	Views           int          `json:"views"`
	Expire          int64        `json:"expire"`
	Preview         bool         `json:"preview"`
	Source          *shareSource `json:"source,omitempty"`
}

// 我的分享列表响应
type ShareList struct {
	Total int           `json:"total"`
	Items []myShareItem `json:"items"`
	User  struct {
		ID    string `json:"id"`
		Nick  string `json:"nick"`
		Group string `json:"group"`
		Date  string `json:"date"`
	} `json:"user"`
}
