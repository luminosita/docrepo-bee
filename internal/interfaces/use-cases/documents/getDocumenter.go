package documents

import (
	"github.com/luminosita/bee/common/interfaces"
	"github.com/luminosita/bee/internal/domain/entities"
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
