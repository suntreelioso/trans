package common

const (
	ApiListFileUrl = "/api/list"
	ApiDownloadUrl = "/api/download"
)

type FileInfo struct {
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type FileList []FileInfo
