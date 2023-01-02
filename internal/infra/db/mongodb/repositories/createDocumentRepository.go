package repositories

import (
	"context"
	"github.com/luminosita/bee/internal/infra/db/mongodb"
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentRepository struct {
	ctx context.Context
	*mongo.Collection
}

func NewDocumentRepository(ctx context.Context) *DocumentRepository {
	return &DocumentRepository{
		ctx: ctx,
	}
}

func (r *DocumentRepository) CreateDocument(docData *documents.RepoRequest) (*documents.RepoResponse, error) {
	_, err := r.getCollection().InsertOne(r.ctx, docData) //, createdAt: new Date());
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *DocumentRepository) getCollection() *mongo.Collection {
	if r.Collection == nil {
		r.Collection = mongodb.GetDbCollection(r.ctx)
	}

	return r.Collection
}
