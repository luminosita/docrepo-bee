package documents

import (
	"github.com/luminosita/honeycomb/pkg/interfaces"
	"github.com/luminosita/sample-bee/internal/domain/entities"
)

type GetDocumenterRequest = struct {
	DocumentId string
}
type GetDocumenterResponse = struct {
	Document *entities.Document
}

type GetDocumenter interface {
	interfaces.UseCaser[*GetDocumenterRequest, *GetDocumenterResponse]
}
