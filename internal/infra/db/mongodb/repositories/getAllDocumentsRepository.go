package repositories

import (
	"context"
	"fmt"
	"github.com/luminosita/bee/internal/infra/db/mongodb"
	"github.com/luminosita/bee/internal/interfaces/respositories/documents"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type GetAllDocumentsRepository struct {
	ctx context.Context
}

func NewGetAllDocumentsRepository(ctx context.Context) *GetAllDocumentsRepository {
	return &GetAllDocumentsRepository{
		ctx: ctx,
	}
}

func (r *GetAllDocumentsRepository) GetAllDocuments(
	docData *documents.GetAllDocumentsRepositorerRequest) (*documents.GetAllDocumentsRepositorerResponse, error) {
	col := mongodb.GetDbCollection(r.ctx)
	cursor, err := col.Find(r.ctx, docData)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all returned documents and print them out.
	//See the mongo.Cursor documentation for more examples of using cursors.
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}

	return nil, nil
}
