package documents

import (
	"github.com/luminosita/sample-bee/internal/interfaces/respositories/documents"
	documents2 "github.com/luminosita/sample-bee/internal/interfaces/use-cases/documents"
)

type CreateDocument struct {
	repo documents.CreateDocumentRepositorer
}

func NewCreateDocument(r documents.CreateDocumentRepositorer) documents2.CreateDocumenter {
	return &CreateDocument{
		repo: r,
	}
}

func (d *CreateDocument) Execute(
	documentData *documents2.CreateDocumenterRequest) (*documents2.CreateDocumenterResponse, error) {
	data := &documents.CreateDocumentRepositorerRequest{
		Document: documentData.Document,
	}

	res, err := d.repo.CreateDocument(data)

	if err != nil {
		return nil, err
	}

	return &documents2.CreateDocumenterResponse{
		DocumentId: res.DocumentId,
	}, nil
}
