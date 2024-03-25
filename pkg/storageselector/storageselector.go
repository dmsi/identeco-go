package storageselector

import (
	"fmt"
	"log/slog"

	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/storage"
	dynamodbdriver "github.com/dmsi/identeco-go/pkg/storage/drivers/dynamodb"
	mongodbdriver "github.com/dmsi/identeco-go/pkg/storage/drivers/mongodb"
	s3driver "github.com/dmsi/identeco-go/pkg/storage/drivers/s3"
)

func NewKeysStorage(lg *slog.Logger) (storage.KeysStorage, error) {
	if _, ok := config.Cfg.KeysStorageDriver.(config.DynamoDBDriverConfig); ok {
		return nil, fmt.Errorf("unable to create keys storage: dynamodb driver not supported")
	} else if _, ok := config.Cfg.KeysStorageDriver.(config.MongoDBDriverConfig); ok {
		return mongodbdriver.NewKeysStorage(lg)
	} else if _, ok := config.Cfg.KeysStorageDriver.(config.S3DriverConfig); ok {
		return s3driver.NewKeysStorage(lg)
	}

	return nil, fmt.Errorf("unable to create keys storage: driver not found")
}

func NewUsersStorage(lg *slog.Logger) (storage.UsersStorage, error) {
	if _, ok := config.Cfg.UserStorageDriver.(config.DynamoDBDriverConfig); ok {
		return dynamodbdriver.NewUsersStorage(lg)
	} else if _, ok := config.Cfg.UserStorageDriver.(config.MongoDBDriverConfig); ok {
		return mongodbdriver.NewUsersStorage(lg)
	} else if _, ok := config.Cfg.UserStorageDriver.(config.S3DriverConfig); ok {
		return nil, fmt.Errorf("unable to create users storage: s3 driver not supported")
	}

	return nil, fmt.Errorf("unable to create keys storage: driver not found")
}
