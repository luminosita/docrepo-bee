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

var GetWireSet = wire.NewSet(NewGetDocumentRepository,
	wire.Bind(new(documents.GetDocumentRepositorer), new(*GetDocumentRepository)))

type GetDocumentRepository struct {
	ctx context.Context
}

func NewGetDocumentRepository(ctx context.Context) *GetDocumentRepository {
	return &GetDocumentRepository{
		ctx: ctx,
	}
}

func (r *GetDocumentRepository) GetDocument(
	docData *documents.GetDocumentRepositorerRequest) (*documents.GetDocumentRepositorerResponse, error) {

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

	//Close() method for stream is not called. Handler is expected to call Close() since
	//DownloadStream implementes Closer interface
	dStream, err := bucket.OpenDownloadStream(docId)
	if err != nil {
		return nil, err
	}

	file := dStream.GetFile()

	return &documents.GetDocumentRepositorerResponse{
		Name:   file.Name,
		Size:   file.Length,
		Reader: dStream,
	}, nil
}
