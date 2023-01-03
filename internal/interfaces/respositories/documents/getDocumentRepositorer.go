package documents

import "github.com/luminosita/sample-bee/internal/domain/entities"

type GetDocumentRepositorerRequest = struct {
	DocumentID string
}
type GetDocumentRepositorerResponse = struct {
	Document *entities.Document
}

type GetDocumentRepositorer interface {
	GetDocument(req *GetDocumentRepositorerRequest) (*GetDocumentRepositorerResponse, error)
}
