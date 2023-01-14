package documents

import (
	"context"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/docrepo-bee/internal/data/db/mongodb"
	"github.com/luminosita/docrepo-bee/internal/interface/data/documents"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GetDocumentInfoRepository struct {
	ctx  context.Context
	data *mongodb.Data
}

func NewGetDocumentInfoRepository(ctx context.Context, data *mongodb.Data) *GetDocumentInfoRepository {
	return &GetDocumentInfoRepository{
		ctx:  ctx,
		data: data,
	}
}

func (r *GetDocumentInfoRepository) GetDocumentInfo(
	docData *documents.GetDocumentInfoRepositorerRequest) (*documents.GetDocumentInfoRepositorerResponse, error) {

	bucket := r.data.DbBucket(Documents)
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
		return nil, log.LogErrorf("cannot find document by ID: %s, %s", docId, err)
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
