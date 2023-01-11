package repositories

import (
	"context"
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	"github.com/luminosita/honeycomb/pkg/infra/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var GetDocumentInfoWireSet = wire.NewSet(NewGetDocumentInfoRepository,
	wire.Bind(new(documents.GetDocumentInfoRepositorer), new(*GetDocumentInfoRepository)))

type GetDocumentInfoRepository struct {
	ctx context.Context
}

func NewGetDocumentInfoRepository(ctx context.Context) *GetDocumentInfoRepository {
	return &GetDocumentInfoRepository{
		ctx: ctx,
	}
}

func (r *GetDocumentInfoRepository) GetDocumentInfo(
	docData *documents.GetDocumentInfoRepositorerRequest) (*documents.GetDocumentInfoRepositorerResponse, error) {

	bucket := mongodb.GetDbBucket(DOCUMENTS)
	fsFiles := bucket.GetFilesCollection()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	docId, err := primitive.ObjectIDFromHex(docData.DocumentId)
	if err != nil {
		return nil, err
	}

	//check if document with specified documentId exists
	var results bson.M
	err = fsFiles.FindOne(ctx, bson.M{"_id": docId}).Decode(&results)
	if err != nil {
		return nil, err
	}

	name := results["filename"].(string)
	size := results["length"].(int64)
	uploadDate := results["uploadDate"].(primitive.DateTime).Time()

	return &documents.GetDocumentInfoRepositorerResponse{
		Name:       name,
		Size:       size,
		UploadDate: uploadDate,
	}, nil
}
