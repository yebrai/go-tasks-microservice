package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMongoConnection = errors.New("error connecting to MongoDB")
)

func GetMongoDBFromURI(ctx context.Context, uri, dbName string) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(uri)
	cli, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMongoConnection, err)
	}
	return cli.Database(dbName), nil
}
