package documents

import (
	"io"
)

type PutDocumentRepositorerRequest = struct {
	Name   string
	Size   int64
	Reader io.Reader
}
type PutDocumentRepositorerResponse = struct {
	DocumentId string
}

type PutDocumentRepositorer interface {
	PutDocument(req *PutDocumentRepositorerRequest) (*PutDocumentRepositorerResponse, error)
}
