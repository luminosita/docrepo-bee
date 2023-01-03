package documents

import (
	"github.com/luminosita/honeycomb/pkg/interfaces"
	"github.com/luminosita/sample-bee/internal/domain/entities"
)

type CreateDocumenterRequest = struct {
	Document *entities.Document
}
type CreateDocumenterResponse = struct {
	DocumentId string
}

type CreateDocumenter interface {
	interfaces.UseCaser[*CreateDocumenterRequest, *CreateDocumenterResponse]
}
