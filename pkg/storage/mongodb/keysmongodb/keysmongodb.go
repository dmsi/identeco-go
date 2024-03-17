package keysmongodb

import (
	"context"
	"log/slog"

	e "github.com/dmsi/identeco-go/pkg/lib/err"
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

type Kkk struct {
	KeyType string `bson:"keytype"`
	KeyData []byte `bson:"keydata"`
}

// index - keytype (private, jwks)
// value - string(bytes)
type KeysStorage struct {
	lg             *slog.Logger
	client         *mongo.Client
	mdb            *mongo.Collection
	databaseName   string
	collectionName string
	ctx            context.Context
}

func wrap(name string, err error) error {
	return e.Wrap("storage.monbodb.keysmongodb."+name, err)
}

func New(lg *slog.Logger, url, database, collection string) (*KeysStorage, error) {
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
			Keys:    bson.D{{Key: indexName, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		theLg.Error("mongodb index failed")
		return nil, wrap("New", err)
	}

	theLg.Info("mongodb index set", "index", index)

	return &KeysStorage{
		lg:             lg,
		client:         client,
		mdb:            mdb,
		databaseName:   database,
		collectionName: collection,
		ctx:            ctx,
	}, nil
}

func (k *KeysStorage) read(keytype string) ([]byte, error) {
	lg := k.lg.With("keytype", keytype)

	filter := bson.D{primitive.E{Key: indexName, Value: keytype}}
	res := k.mdb.FindOne(k.ctx, filter)

	key := struct {
		KeyType string `bson:"keytype"`
		KeyData []byte `bson:"keydata"`
	}{}
	err := res.Decode(&key)
	if err != nil {
		lg.Error("read failed")
		return nil, wrap("read", err)
	}

	return key.KeyData, nil
}

func (k *KeysStorage) write(keytype string, data []byte) error {
	lg := k.lg.With("keytype", keytype)

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
		lg.Error("write failed")
		return wrap("write", err)
	}

	return nil
}

func (k *KeysStorage) ReadPrivateKey() (*storage.PrivateKeyData, error) {
	data, err := k.read(privateKeyType)
	if err != nil {
		return nil, wrap("ReadPrivateKey", err)
	}

	return &storage.PrivateKeyData{
		Data: data,
	}, nil
}

func (k *KeysStorage) WritePrivateKey(key storage.PrivateKeyData) error {
	err := k.write(privateKeyType, key.Data)
	if err != nil {
		return wrap("WritePrivateKey", err)
	}

	return nil
}

func (k *KeysStorage) ReadJWKSets() (*storage.JWKSetsData, error) {
	data, err := k.read(jwkSetsType)
	if err != nil {
		return nil, wrap("ReadJWKSets", err)
	}

	return &storage.JWKSetsData{
		Data: data,
	}, nil
}

func (k *KeysStorage) WriteJWKSets(jwkSets storage.JWKSetsData) error {
	err := k.write(jwkSetsType, jwkSets.Data)
	if err != nil {
		return wrap("WriteJWKSets", err)
	}

	return nil
}
