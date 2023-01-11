package repositories

import (
	"context"
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	"github.com/luminosita/honeycomb/pkg/infra/db/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var PutDocumentWireSet = wire.NewSet(NewPutDocumentRepository,
	wire.Bind(new(documents.PutDocumentRepositorer), new(*PutDocumentRepository)))

type PutDocumentRepository struct {
	ctx context.Context
}

func NewPutDocumentRepository(ctx context.Context) *PutDocumentRepository {
	return &PutDocumentRepository{
		ctx: ctx,
	}
}

func (r *PutDocumentRepository) PutDocument(
	docData *documents.PutDocumentRepositorerRequest) (res *documents.PutDocumentRepositorerResponse, err error) {

	bucket := mongodb.GetDbBucket(DOCUMENTS)

	docId := primitive.NewObjectID()

	uploadStream, err := bucket.OpenUploadStreamWithID(docId, docData.Name)
	if err != nil {
		return nil, err
	}

	return &documents.PutDocumentRepositorerResponse{
		DocumentId: docId.Hex(),
		Writer:     uploadStream,
	}, nil
}
