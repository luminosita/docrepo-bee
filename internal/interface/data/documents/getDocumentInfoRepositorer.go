//go:generate mockgen -destination=./mocks/mock_getDocumentInfoRepositorer.go -package=mocks . GetDocumentInfoRepositorer
package documents

import "time"

type GetDocumentInfoRepositorerRequest struct {
	DocumentId string
}

type GetDocumentInfoRepositorerResponse struct {
	Name       string
	Size       int64
	UploadDate time.Time
}

type GetDocumentInfoRepositorer interface {
	GetDocumentInfo(req *GetDocumentInfoRepositorerRequest) (*GetDocumentInfoRepositorerResponse, error)
}
