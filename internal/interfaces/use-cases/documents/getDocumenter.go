//go:generate mockgen -destination=./mocks/mock_getDocumenter.go -package=mocks . GetDocumenter
package documents

import (
	"github.com/luminosita/honeycomb/pkg/interfaces"
	"io"
)

type GetDocumenterRequest struct {
	DocumentId string
}

type GetDocumenterResponse struct {
	Name   string
	Size   int64
	Reader io.Reader
}

type GetDocumenter interface {
	interfaces.UseCaser[*GetDocumenterRequest, *GetDocumenterResponse]
}
