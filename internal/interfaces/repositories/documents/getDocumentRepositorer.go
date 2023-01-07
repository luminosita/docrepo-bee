//go:generate mockgen -destination=./mocks/mock_getDocumentRepositorer.go -package=mocks . GetDocumentRepositorer
package documents

import (
	"io"
)

type GetDocumentRepositorerRequest struct {
	DocumentId string
}

type GetDocumentRepositorerResponse struct {
	Name   string
	Size   int64
	Reader io.Reader
}

type GetDocumentRepositorer interface {
	GetDocument(req *GetDocumentRepositorerRequest) (*GetDocumentRepositorerResponse, error)
}
