package repositories

import (
	"context"
	"github.com/luminosita/bee/internal/domain/entities"
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
)

type MockDocumentsRepository struct {
	ctx context.Context
}

func MakeMockDocumentsRepository(ctx context.Context) *MockDocumentsRepository {
	return &MockDocumentsRepository{
		ctx: ctx,
	}
}

func (r *MockDocumentsRepository) GetAllDocuments(
	docData *documents.GetAllDocumentsRepositorerRequest) (*documents.GetAllDocumentsRepositorerResponse, error) {

	return &documents.GetAllDocumentsRepositorerResponse{
		Documents: append(make([]*entities.Document, 0), &entities.Document{
			Name:       "MockName",
			DocumentId: "MockId",
		}),
	}, nil
}

func (r *MockDocumentsRepository) GetDocument(
	docData *documents.GetDocumentRepositorerRequest) (*documents.GetDocumentRepositorerResponse, error) {
	return &documents.GetDocumentRepositorerResponse{
		Document: &entities.Document{
			Name:       "MockName",
			DocumentId: "MockId",
		},
	}, nil
}

func (r *MockDocumentsRepository) CreateDocument(
	docData *documents.CreateDocumentRepositorerRequest) (*documents.CreateDocumentRepositorerResponse, error) {
	return &documents.CreateDocumentRepositorerResponse{
		DocumentId: "MockId",
	}, nil
}
