package serializer

import (
	"time"
)

// User 用户序列化器
type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"user_name"`
	Nickname       string    `json:"nickname"`
	Status         int       `json:"status"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	PreferredTheme string    `json:"preferred_theme"`
	Score          int       `json:"score"`
	Anonymous      bool      `json:"anonymous"`
	Group          group     `json:"group"`
	Tags           []tag     `json:"tags"`
}

type group struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	AllowShare           bool   `json:"allowShare"`
	AllowRemoteDownload  bool   `json:"allowRemoteDownload"`
	AllowArchiveDownload bool   `json:"allowArchiveDownload"`
	ShareFreeEnabled     bool   `json:"shareFree"`
	ShareDownload        bool   `json:"shareDownload"`
	CompressEnabled      bool   `json:"compress"`
	WebDAVEnabled        bool   `json:"webdav"`
	RelocateEnabled      bool   `json:"relocate"`
	SourceBatchSize      int    `json:"sourceBatch"`
	SelectNode           bool   `json:"selectNode"`
	AdvanceDelete        bool   `json:"advanceDelete"`
	AllowWebDAVProxy     bool   `json:"allowWebDAVProxy"`
}

type tag struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Color      string `json:"color"`
	Type       int    `json:"type"`
	Expression string `json:"expression"`
}

type Storage struct {
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Total uint64 `json:"total"`
}

// WebAuthnCredentials 外部验证器凭证
type WebAuthnCredentials struct {
	ID          []byte `json:"id"`
	FingerPrint string `json:"fingerprint"`
}
