package documents

import (
	"fmt"
	"github.com/luminosita/common-bee/pkg/log"
	documents2 "github.com/luminosita/docrepo-bee/internal/interface/app/documents"
	"github.com/luminosita/docrepo-bee/internal/interface/data/documents"
)

type PutDocument struct {
	repo documents.PutDocumentRepositorer
}

func NewPutDocument(r documents.PutDocumentRepositorer) *PutDocument {
	return &PutDocument{
		repo: r,
	}
}

func (d *PutDocument) Execute(
	docData *documents2.PutDocumenterRequest) (*documents2.PutDocumenterResponse, error) {

	if docData == nil || len(docData.Name) == 0 || docData.Size <= 0 {
		return nil, log.LogErrorf(fmt.Sprintf("Bad request: %+v", docData))
	}

	data := &documents.PutDocumentRepositorerRequest{
		Name: docData.Name,
		Size: docData.Size,
	}

	res, err := d.repo.PutDocument(data)
	if err != nil {
		return nil, err
	}

	return &documents2.PutDocumenterResponse{
		DocumentId: res.DocumentId,
		Writer:     res.Writer,
	}, nil
}
