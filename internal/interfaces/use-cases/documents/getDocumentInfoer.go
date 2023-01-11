//go:generate mockgen -destination=./mocks/mock_getDocumentInfoer.go -package=mocks . GetDocumentInfoer
package documents

import (
	"github.com/luminosita/honeycomb/pkg/interfaces"
	"time"
)

type GetDocumentInfoerRequest struct {
	DocumentId string
}

type GetDocumentInfoerResponse struct {
	Name       string
	Size       int64
	UploadDate time.Time
}

type GetDocumentInfoer interface {
	interfaces.UseCaser[*GetDocumentInfoerRequest, *GetDocumentInfoerResponse]
}
