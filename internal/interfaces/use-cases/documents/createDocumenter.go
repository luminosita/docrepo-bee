package documents

import "github.com/luminosita/bee/common/interfaces"

type CreateDocumenterRequest = struct {
	Content string
}
type CreateDocumenterResponse = struct {
	DocumentId string
}

type CreateDocumenter interface {
	interfaces.UseCaser[*CreateDocumenterRequest, *CreateDocumenterResponse]
}
