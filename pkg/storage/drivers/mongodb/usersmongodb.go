package mongodbdriver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersMongoDBDriver struct {
	lg             *slog.Logger
	client         *mongo.Client
	mdb            *mongo.Collection
	databaseName   string
	collectionName string
	ctx            context.Context
}

// TODO what if both keydata and userdata want to share the same connection?
// storage.Session {} interface with open/close?
func NewUsersStorage(lg *slog.Logger) (*UsersMongoDBDriver, error) {
	cfg, ok := config.Cfg.UserStorageDriver.(config.MongoDBDriverConfig)
	if !ok {
		return nil, fmt.Errorf("users driver configuration is not provided: mongodb")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ConnectURL))
	if err != nil {
		return nil, err
	}

	mdb := client.Database(cfg.Database).Collection(cfg.Collection)

	// set username as unique index
	_, err = mdb.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	return &UsersMongoDBDriver{
		lg:             lg,
		client:         client,
		mdb:            mdb,
		databaseName:   cfg.Database,
		collectionName: cfg.Collection,
		ctx:            ctx,
	}, nil
}

func (s UsersMongoDBDriver) ReadUserData(username string) (*storage.UserData, error) {
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	res := s.mdb.FindOne(s.ctx, filter)

	mongoUser := struct {
		Username string `bson:"username"`
		Hash     string `bson:"hash"`
	}{}
	err := res.Decode(&mongoUser)
	if err != nil {
		return nil, err
	}

	return &storage.UserData{
		Username: mongoUser.Username,
		Hash:     mongoUser.Hash,
	}, nil
}

func (s UsersMongoDBDriver) WriteUserData(user storage.UserData) error {
	if user.Username == "" || user.Hash == "" {
		return errors.New("invalid arguments")
	}

	mongoUser := struct {
		Username string `bson:"username"`
		Hash     string `bson:"hash"`
	}{
		Username: user.Username,
		Hash:     user.Hash,
	}

	_, err := s.mdb.InsertOne(s.ctx, &mongoUser)
	if err != nil {
		return err
	}

	return nil
}
