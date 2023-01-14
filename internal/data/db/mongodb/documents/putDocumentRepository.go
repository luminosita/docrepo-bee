package documents

import (
	"context"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/docrepo-bee/internal/data/db/mongodb"
	"github.com/luminosita/docrepo-bee/internal/interface/data/documents"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PutDocumentRepository struct {
	ctx    context.Context
	data   *mongodb.Data
	logger log.LogHelper
}

func NewPutDocumentRepository(ctx context.Context, data *mongodb.Data) *PutDocumentRepository {
	return &PutDocumentRepository{
		ctx:  ctx,
		data: data,
	}
}

func (r *PutDocumentRepository) PutDocument(
	docData *documents.PutDocumentRepositorerRequest) (res *documents.PutDocumentRepositorerResponse, err error) {

	bucket := r.data.DbBucket(Documents)

	docId := primitive.NewObjectID()

	uploadStream, err := bucket.OpenUploadStreamWithID(docId, docData.Name)
	if err != nil {
		return nil, log.LogErrorf("cannot open bucket upload stream: %s", err)
	}

	return &documents.PutDocumentRepositorerResponse{
		DocumentId: docId.Hex(),
		Writer:     uploadStream,
	}, nil
}
