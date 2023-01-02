package repositories

import (
	"context"
	"github.com/luminosita/bee/internal/infra/db/mongodb"
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetDocumentRepository struct {
	ctx context.Context
	*mongo.Collection
}

func NewGetDocumentRepository(ctx context.Context) *GetDocumentRepository {
	return &GetDocumentRepository{
		ctx: ctx,
	}
}

func (r *GetDocumentRepository) GetDocument(
	docData *documents.GetDocumentRepositorerRequest) (*documents.GetDocumentRepositorerResponse, error) {
	_ = r.getCollection().FindOne(r.ctx, docData) //, createdAt: new Date());

	return nil, nil
}

func (r *GetDocumentRepository) getCollection() *mongo.Collection {
	if r.Collection == nil {
		r.Collection = mongodb.GetDbCollection(r.ctx)
	}

	return r.Collection
}
