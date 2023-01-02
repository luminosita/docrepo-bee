package documents

import (
	"github.com/luminosita/bee/common/interfaces"
)

type GetAllDocumenterRequest = struct {
	DocumentId string
}
type GetAllDocumenterResponse = struct {
	Content string
}

type GetAllDocumenter interface {
	interfaces.UseCaser[*GetAllDocumenterRequest, *GetAllDocumenterResponse]
}
