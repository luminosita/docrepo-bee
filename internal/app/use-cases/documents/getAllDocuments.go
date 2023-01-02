package documents

import (
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
	documents2 "github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

type GetDocument struct {
	repo documents.GetDocumentRepositorer
}

func NewGetDocument(r documents.GetDocumentRepositorer) documents2.GetDocumenter {
	return &GetDocument{
		repo: r,
	}
}

func (d *GetDocument) Execute(
	documentData *documents2.GetDocumenterRequest) (*documents2.GetDocumenterResponse, error) {
	data := &documents.GetDocumentRepositorerRequest{}

	_, err := d.repo.GetDocument(data)

	if err != nil {
		return nil, err
	}

	//model conversion

	return nil, nil
}
