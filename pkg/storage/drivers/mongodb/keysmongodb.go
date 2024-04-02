package mongodbdriver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	indexName      = "keytype"
	privateKeyType = "private"
	jwkSetsType    = "jwks"
)

// index - keytype (private, jwks)
// value - string(bytes)
type KeysMongoDBDriver struct {
	lg             *slog.Logger
	client         *mongo.Client
	mdb            *mongo.Collection
	databaseName   string
	collectionName string
	ctx            context.Context
}

func NewKeysStorage(lg *slog.Logger) (*KeysMongoDBDriver, error) {
	cfg, ok := config.Cfg.KeysStorageDriver.(config.MongoDBDriverConfig)
	if !ok {
		return nil, fmt.Errorf("keys driver configuration is not provided: mongodb")
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
			Keys:    bson.D{{Key: indexName, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	return &KeysMongoDBDriver{
		lg:             lg,
		client:         client,
		mdb:            mdb,
		databaseName:   cfg.Database,
		collectionName: cfg.Collection,
		ctx:            ctx,
	}, nil
}

func (k *KeysMongoDBDriver) read(keytype string) ([]byte, error) {
	filter := bson.D{primitive.E{Key: indexName, Value: keytype}}
	res := k.mdb.FindOne(k.ctx, filter)

	key := struct {
		KeyType string `bson:"keytype"`
		KeyData []byte `bson:"keydata"`
	}{}
	err := res.Decode(&key)
	if err != nil {
		return nil, err
	}

	return key.KeyData, nil
}

func (k *KeysMongoDBDriver) write(keytype string, data []byte) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: indexName, Value: keytype}}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: indexName, Value: keytype},
			primitive.E{Key: "keydata", Value: data},
		},
	}}

	_, err := k.mdb.UpdateOne(k.ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s KeysMongoDBDriver) ReadKeys() (*storage.Keys, error) {
	privateKey, err := s.read(privateKeyType)
	if err != nil {
		return nil, err
	}

	jwks, err := s.read(jwkSetsType)
	if err != nil {
		return nil, err
	}

	return &storage.Keys{
		PrivateKey: privateKey,
		JWKS:       jwks,
	}, nil
}

func (s KeysMongoDBDriver) WriteKeys(k storage.Keys) error {
	err := s.write(privateKeyType, k.PrivateKey)
	if err != nil {
		return err
	}

	err = s.write(jwkSetsType, k.JWKS)
	if err != nil {
		return err
	}

	return nil
}
