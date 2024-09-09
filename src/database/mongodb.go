package database

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitializeMongoDatabase(ctx context.Context, dsn, dbName string) *mongo.Database {
	log := logrus.WithContext(ctx)

	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("error when mongo.Connect(ctx, clientOptions), request: %v, error: %v", dsn, err.Error())
		return nil
	}

	mongoDB := client.Database(dbName)

	log.Infof("ping mongodb %s", dbName)
	if err := mongoDB.Client().Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("error when ping mongodb %s, error: %v", dbName, err.Error())
	}

	return mongoDB
}
