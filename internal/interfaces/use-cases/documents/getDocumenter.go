package documents

import (
	"github.com/luminosita/bee/common/interfaces"
)

type GetDocumenterRequest = struct {
	DocumentId string
}
type GetDocumenterResponse = struct {
	Content string
}

type GetDocumenter interface {
	interfaces.UseCaser[*GetDocumenterRequest, *GetDocumenterResponse]
}
