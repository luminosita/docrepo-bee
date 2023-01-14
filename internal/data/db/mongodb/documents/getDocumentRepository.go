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

type GetDocumentRepository struct {
	ctx  context.Context
	data *mongodb.Data
}

func NewGetDocumentRepository(ctx context.Context, data *mongodb.Data) *GetDocumentRepository {
	return &GetDocumentRepository{
		ctx:  ctx,
		data: data,
	}
}

func (r *GetDocumentRepository) GetDocument(
	docData *documents.GetDocumentRepositorerRequest) (*documents.GetDocumentRepositorerResponse, error) {

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

	//Close() method for stream is not called. Handler is expected to call Close() since
	//DownloadStream implementes Closer interface
	dStream, err := bucket.OpenDownloadStream(docId)
	if err != nil {
		return nil, log.LogErrorf("cannot open bucket download stream: %s", err)
	}

	file := dStream.GetFile()

	return &documents.GetDocumentRepositorerResponse{
		Name:   file.Name,
		Size:   file.Length,
		Reader: dStream,
	}, nil
}
