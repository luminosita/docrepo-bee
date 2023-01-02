package mongodb

import (
	"context"
	"fmt"
	rkmongo "github.com/rookie-ninja/rk-db/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

const DOCUMENT_COLLECTION string = "documents"

var (
	once sync.Once
	db   *mongo.Database
)

func GetDbCollection(ctx context.Context) *mongo.Collection {
	once.Do(func() { // <-- atomic, does not allow repeating
		db := rkmongo.GetMongoDB("bee-mongo", "bee")
		createCollection(ctx, db, DOCUMENT_COLLECTION)
	})

	return db.Collection(DOCUMENT_COLLECTION)
}

func createCollection(ctx context.Context, db *mongo.Database, name string) {
	opts := options.CreateCollection()
	err := db.CreateCollection(ctx, name, opts)
	if err != nil {
		fmt.Println("collection exists may be, continue")
	}
}
