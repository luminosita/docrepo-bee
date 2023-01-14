//go:generate mockgen -destination=./mocks/mock_getDocumentInfoer.go -package=mocks . GetDocumentInfoer
package documents

import (
	"github.com/luminosita/docrepo-bee/internal/interface/app"
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
	app.UseCaser[*GetDocumentInfoerRequest, *GetDocumentInfoerResponse]
}
