package documents

import "github.com/luminosita/bee/internal/domain/entities"

type CreateDocumentRepositorerRequest = struct {
	Document *entities.Document
}
type CreateDocumentRepositorerResponse = struct {
	DocumentId string
}

type CreateDocumentRepositorer interface {
	CreateDocument(req *CreateDocumentRepositorerRequest) (*CreateDocumentRepositorerResponse, error)
}
