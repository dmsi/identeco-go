package usersmongodb

import (
	"context"
	"errors"

	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

type UsersStorage struct {
	lg             *slog.Logger
	client         *mongo.Client
	mdb            *mongo.Collection
	databaseName   string
	collectionName string
	ctx            context.Context
}

func wrap(name string, err error) error {
	return e.Wrap("storage.mongodb.usersmongodb."+name, err)
}

// TODO what if both keydata and userdata want to share the same connection?
// storage.Session {} interface with open/close?
func New(lg *slog.Logger, url, database, collection string) (*UsersStorage, error) {
	ctx := context.Background()
	theLg := lg.With("database", database, "collection", collection)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, wrap("New", err)
	}

	theLg.Info("mongodb connected")

	mdb := client.Database(database).Collection(collection)

	// set username as unique index
	index, err := mdb.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		theLg.Error("mongodb index failed")
		return nil, wrap("New", err)
	}

	theLg.Info("mongodb index set", "index", index)

	return &UsersStorage{
		lg:             lg,
		client:         client,
		mdb:            mdb,
		databaseName:   database,
		collectionName: collection,
		ctx:            ctx,
	}, nil
}

func (u *UsersStorage) ReadUserData(username string) (*storage.UserData, error) {
	filter := bson.D{{"username", username}}
	res := u.mdb.FindOne(u.ctx, filter)

	mongoUser := struct {
		Username string `bson:"username"`
		Hash     string `bson:"hash"`
	}{}
	err := res.Decode(&mongoUser)
	if err != nil {
		return nil, wrap("ReadUserData", err)
	}

	return &storage.UserData{
		Username: mongoUser.Username,
		Hash:     mongoUser.Hash,
	}, nil
}

func (u *UsersStorage) WriteUserData(user storage.UserData) error {
	if user.Username == "" || user.Hash == "" {
		return wrap("WriteUserData", errors.New("invalid arguments"))
	}

	mongoUser := struct {
		Username string `bson:"username"`
		Hash     string `bson:"hash"`
	}{
		Username: user.Username,
		Hash:     user.Hash,
	}

	_, err := u.mdb.InsertOne(u.ctx, &mongoUser)
	if err != nil {
		return wrap("WriteUserData", err)
	}

	return nil
}
