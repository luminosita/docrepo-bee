package mongodb

import (
	"context"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/docrepo-bee/internal/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Data struct {
	db     *mongo.Database
	col    *mongo.Collection
	bucket *gridfs.Bucket

	c *conf.Data
}

// NewData .
func NewData(ctx context.Context, c *conf.Data) (*Data, func()) {
	db := initDb(ctx, c)

	cleanup := func() {
		log.Info("closing the data resources")
		if err := db.Client().Disconnect(ctx); err != nil {
			panic(err)
		}
	}
	return &Data{c: c, db: db}, cleanup
}

func (d *Data) DbBucket(name string) *gridfs.Bucket {
	if d.bucket == nil {
		d.bucket = createBucket(d.db, name)
	}

	return d.bucket
}

func initDb(ctx context.Context, c *conf.Data) *mongo.Database {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Database.Source))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Info("Successfully connected and pinged.")

	return client.Database(c.Database.Db)
}

func createBucket(db *mongo.Database, name string) *gridfs.Bucket {
	opts := &options.BucketOptions{
		Name: &name,
	}
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}

	return bucket
}
