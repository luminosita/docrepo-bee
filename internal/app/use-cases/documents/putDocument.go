package documents

import (
	"github.com/luminosita/docrepo-bee/internal/interfaces/respositories/documents"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
)

type PutDocument struct {
	repo documents.PutDocumentRepositorer
}

func NewPutDocument(r documents.PutDocumentRepositorer) documents2.PutDocumenter {
	return &PutDocument{
		repo: r,
	}
}

func (d *PutDocument) Execute(
	docData *documents2.PutDocumenterRequest) (*documents2.PutDocumenterResponse, error) {
	data := &documents.PutDocumentRepositorerRequest{
		Name:   docData.Name,
		Size:   docData.Size,
		Reader: docData.Reader,
	}

	res, err := d.repo.PutDocument(data)
	if err != nil {
		return nil, err
	}

	return &documents2.PutDocumenterResponse{
		DocumentId: res.DocumentId,
	}, nil
}
