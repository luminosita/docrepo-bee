package repositories

import (
	"bufio"
	"context"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	"github.com/luminosita/honeycomb/pkg/infra/db/mongodb"
	"github.com/luminosita/honeycomb/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

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
	defer func() {
		err = uploadStream.Close()
	}()

	reader := bufio.NewReader(docData.Reader)

	buf := make([]byte, 4*1024) //the chunk size

	bytesWritten, err := io.CopyBuffer(uploadStream, reader, buf)
	if err != nil {
		return nil, err
	}

	log.Log().Infof("Write file to Bucket was successful. File size: %d B\n", bytesWritten)

	return &documents.PutDocumentRepositorerResponse{
		DocumentId: docId.Hex(),
	}, nil
}
