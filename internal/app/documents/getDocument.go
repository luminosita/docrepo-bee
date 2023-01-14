package documents

import (
	"fmt"
	"github.com/luminosita/common-bee/pkg/log"
	documents2 "github.com/luminosita/docrepo-bee/internal/interface/app/documents"
	"github.com/luminosita/docrepo-bee/internal/interface/data/documents"
)

type GetDocument struct {
	repo documents.GetDocumentRepositorer
}

func NewGetDocument(r documents.GetDocumentRepositorer) *GetDocument {
	return &GetDocument{
		repo: r,
	}
}

func (d *GetDocument) Execute(
	docData *documents2.GetDocumenterRequest) (*documents2.GetDocumenterResponse, error) {

	if docData == nil || len(docData.DocumentId) == 0 {
		return nil, log.LogErrorf(fmt.Sprintf("Bad request: %+v", docData))
	}

	data := &documents.GetDocumentRepositorerRequest{
		DocumentId: docData.DocumentId,
	}

	res, err := d.repo.GetDocument(data)

	if err != nil {
		return nil, err
	}

	return &documents2.GetDocumenterResponse{
		Name:   res.Name,
		Size:   res.Size,
		Reader: res.Reader,
	}, nil
}
