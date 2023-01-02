package documents

import (
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
	documents2 "github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

type CreateDocument struct {
	repo documents.CreateDocumentRepositorer
}

func NewCreateDocument(r documents.CreateDocumentRepositorer) documents2.CreateDocumenter {
	return &CreateDocument{
		repo: r,
	}
}

func (d *CreateDocument) Execute(documentData *documents2.Request) (*documents2.Response, error) {
	data := &documents.RepoRequest{}

	_, err := d.repo.CreateDocument(data)

	if err != nil {
		return nil, err
	}

	//model conversion

	return nil, nil
}
