package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	authOpt = options.Credential{
		AuthMechanism: os.Getenv("AUTH"),
		AuthSource:    os.Getenv("DB_NAME"),
		Username:      os.Getenv("DB_USER"),
		Password:      os.Getenv("DB_PWD"),
	}
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
)

func ConnectDB(ctx context.Context, host string, port string) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", host, port)
	clientopt := options.Client().SetAuth(authOpt).ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientopt)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func Write(ctx context.Context, dbName string, collName string, args interface{}) error {
	wctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := ConnectDB(wctx, host, port)

	if err != nil {
		return err
	}

	defer client.Disconnect(wctx)

	coll := client.Database(dbName).Collection(collName)

	if _, err := coll.InsertOne(wctx, args); err != nil {
		return err
	}

	return nil
}
