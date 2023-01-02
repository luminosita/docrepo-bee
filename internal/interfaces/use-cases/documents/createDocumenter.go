package documents

import (
	"github.com/luminosita/bee/common/interfaces"
	"github.com/luminosita/bee/internal/domain/entities"
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
