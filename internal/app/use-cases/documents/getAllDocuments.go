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

	res, err := d.repo.GetAllDocuments(data)

	if err != nil {
		return nil, err
	}

	return &documents2.GetAllDocumenterResponse{
		Documents: res.Documents,
	}, nil
}
