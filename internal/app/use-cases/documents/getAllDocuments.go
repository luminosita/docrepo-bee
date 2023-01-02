package documents

import (
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
	documents2 "github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

type GetAllDocuments struct {
	repo documents.GetAllDocumentsRepositorer
}

func NewGetAllDocuments(r documents.GetAllDocumentsRepositorer) documents2.GetAllDocumenter {
	return &GetAllDocuments{
		repo: r,
	}
}

func (d *GetAllDocuments) Execute(
	documentData *documents2.GetAllDocumenterRequest) (*documents2.GetAllDocumenterResponse, error) {
	data := &documents.GetAllDocumentsRepositorerRequest{}

	_, err := d.repo.GetAllDocuments(data)

	if err != nil {
		return nil, err
	}

	//model conversion

	return nil, nil
}
