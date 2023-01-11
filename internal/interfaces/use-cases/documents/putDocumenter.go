//go:generate mockgen -destination=./mocks/mock_putDocumenter.go -package=mocks . PutDocumenter
package documents

import (
	"github.com/luminosita/honeycomb/pkg/interfaces"
	"io"
)

type PutDocumenterRequest struct {
	Name string
	Size int64
}
type PutDocumenterResponse struct {
	DocumentId string
	Writer     io.WriteCloser
}

type PutDocumenter interface {
	interfaces.UseCaser[*PutDocumenterRequest, *PutDocumenterResponse]
}
