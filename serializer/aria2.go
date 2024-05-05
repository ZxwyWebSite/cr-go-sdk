package serializer

import (
	"time"

	"github.com/ZxwyWebSite/cr-go-sdk/pkg/aria2/rpc"
)

// DownloadListResponse 下载列表响应条目
type DownloadListResponse struct {
	UpdateTime     time.Time      `json:"update"`
	UpdateInterval int            `json:"interval"`
	Name           string         `json:"name"`
	Status         int            `json:"status"`
	Dst            string         `json:"dst"`
	Total          uint64         `json:"total"`
	Downloaded     uint64         `json:"downloaded"`
	Speed          int            `json:"speed"`
	Info           rpc.StatusInfo `json:"info"`
	NodeName       string         `json:"node"`
}

// FinishedListResponse 已完成任务条目
type FinishedListResponse struct {
	Name       string         `json:"name"`
	GID        string         `json:"gid"`
	Status     int            `json:"status"`
	Dst        string         `json:"dst"`
	Error      string         `json:"error"`
	Total      uint64         `json:"total"`
	Files      []rpc.FileInfo `json:"files"`
	TaskStatus int            `json:"task_status"`
	TaskError  string         `json:"task_error"`
	CreateTime time.Time      `json:"create"`
	UpdateTime time.Time      `json:"update"`
	NodeName   string         `json:"node"`
}
