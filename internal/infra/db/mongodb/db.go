package config

import (
	"context"
	"fmt"
	rkmongo "github.com/rookie-ninja/rk-db/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: DB connection singelton
var (
	// GlobalAppCtx global application context
	GlobalBeeCtx = &beeContext{}
)

type beeContext struct {
	documentCollection *mongo.Collection
}

func createCollection(ctx context.Context, db *mongo.Database, name string) {
	opts := options.CreateCollection()
	err := db.CreateCollection(ctx, name, opts)
	if err != nil {
		fmt.Println("collection exists may be, continue")
	}
}

func (bc *beeContext) GetDbCollection(ctx context.Context) *mongo.Collection {
	if bc.documentCollection == nil {
		db := rkmongo.GetMongoDB("bee-mongo", "bee")
		createCollection(ctx, db, "meta")
		bc.documentCollection = db.Collection("documents")
	}

	return bc.documentCollection
}
