//go:generate mockgen -destination=./mocks/mock_getDocumenter.go -package=mocks . GetDocumenter
package documents

import (
	"github.com/luminosita/docrepo-bee/internal/interface/app"
	"io"
)

type GetDocumenterRequest struct {
	DocumentId string
}

type GetDocumenterResponse struct {
	Name   string
	Size   int64
	Reader io.ReadCloser
}

type GetDocumenter interface {
	app.UseCaser[*GetDocumenterRequest, *GetDocumenterResponse]
}
