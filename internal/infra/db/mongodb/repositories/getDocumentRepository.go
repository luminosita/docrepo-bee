package repositories

import (
	"context"
	"github.com/luminosita/bee/internal/infra/db/mongodb"
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateDocumentRepository struct {
	ctx context.Context
	*mongo.Collection
}

func NewCreateDocumentRepository(ctx context.Context) *CreateDocumentRepository {
	return &CreateDocumentRepository{
		ctx: ctx,
	}
}

func (r *CreateDocumentRepository) CreateDocument(
	docData *documents.CreateDocumentRepositorerRequest) (*documents.CreateDocumentRepositorerResponse, error) {
	_, err := r.getCollection().InsertOne(r.ctx, docData) //, createdAt: new Date());
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *CreateDocumentRepository) getCollection() *mongo.Collection {
	if r.Collection == nil {
		r.Collection = mongodb.GetDbCollection(r.ctx)
	}

	return r.Collection
}
