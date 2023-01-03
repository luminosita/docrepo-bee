package repositories

import (
	"context"
	"github.com/luminosita/sample-bee/internal/infra/db/mongodb"
	"github.com/luminosita/sample-bee/internal/interfaces/respositories/documents"
)

type CreateDocumentRepository struct {
	ctx context.Context
}

func NewCreateDocumentRepository(ctx context.Context) *CreateDocumentRepository {
	return &CreateDocumentRepository{
		ctx: ctx,
	}
}

func (r *CreateDocumentRepository) CreateDocument(
	docData *documents.CreateDocumentRepositorerRequest) (*documents.CreateDocumentRepositorerResponse, error) {
	col := mongodb.GetDbCollection(r.ctx)
	_, err := col.InsertOne(r.ctx, docData) //, createdAt: new Date());
	if err != nil {
		return nil, err
	}

	return nil, nil
}
