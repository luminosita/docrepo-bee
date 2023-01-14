//go:generate mockgen -destination=./mocks/mock_putDocumentRepositorer.go -package=mocks . PutDocumentRepositorer
package documents

import (
	"io"
)

type PutDocumentRepositorerRequest struct {
	Name string
	Size int64
}
type PutDocumentRepositorerResponse struct {
	DocumentId string
	Writer     io.WriteCloser
}

type PutDocumentRepositorer interface {
	PutDocument(req *PutDocumentRepositorerRequest) (*PutDocumentRepositorerResponse, error)
}
