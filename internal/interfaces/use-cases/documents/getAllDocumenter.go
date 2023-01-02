package documents

import (
	"github.com/luminosita/bee/common/interfaces"
	"github.com/luminosita/bee/internal/domain/entities"
)

type GetAllDocumenterRequest = struct {
}
type GetAllDocumenterResponse = struct {
	Documents []*entities.Document
}

type GetAllDocumenter interface {
	interfaces.UseCaser[*GetAllDocumenterRequest, *GetAllDocumenterResponse]
}
