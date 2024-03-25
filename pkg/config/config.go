package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

const (
	envStage                = "IDO_DEPLOYMENT_STAGE"
	envPrivateKeyLength     = "IDO_PRIVATE_KEY_LENGTH"
	envAccessTokenLifetime  = "IDO_ACCESS_TOKEN_LIFETIME"
	envRefreshTokenLifetime = "IDO_REFRESH_TOKEN_LIFETIME"
	envIssClaim             = "IDO_CLAIM_ISS"
	envTableName            = "IDO_TABLE_NAME"
	envBucketName           = "IDO_BUCKET_NAME"
	envPrivateKeyName       = "IDO_PRIVATE_KEY_NAME"
	envJWKSetsName          = "IDO_JWKS_NAME"
	envStorageDriverPrefix  = "IDO_STORAGE_DRIVER"
	envStorageDriverUsers   = "IDO_STORAGE_DRIVER_USERS"
	envStorageDriverKeys    = "IDO_STORAGE_DRIVER_KEYS"
)

type MongoDBDriverConfig struct {
	ConnectURL string
	Database   string
	Collection string
}

type DynamoDBDriverConfig struct {
	TableName string
	// PrimaryKey string
}

type S3DriverConfig struct {
	BucketName     string
	PrivateKeyName string
	JWKSName       string
}

type Config struct {
	TokenIssClaim        string
	TokenRefreshDuration time.Duration
	TokenAccessDuration  time.Duration
	KeyLength            int
	Stage                string
	UserStorageDriver    any
	KeysStorageDriver    any
}

var Cfg Config

func init() {
	err := readConfig()
	if err != nil {
		log.Fatalf("Unable to read config: %v", err)
	}
}

func readMongoDBConfig(storage string) (*MongoDBDriverConfig, error) {
	urlKey := fmt.Sprintf("%s_%s_MONGODB_URL", envStorageDriverPrefix, strings.ToUpper(storage))
	dbKey := fmt.Sprintf("%s_%s_MONGODB_DATABASE", envStorageDriverPrefix, strings.ToUpper(storage))
	collectionKey := fmt.Sprintf("%s_%s_MONGODB_COLLECTION", envStorageDriverPrefix, strings.ToUpper(storage))

	url := readString(urlKey)
	if url == nil {
		return nil, fmt.Errorf("validation error: %s MongoDB url not set", storage)
	}

	db := readString(dbKey)
	if db == nil {
		return nil, fmt.Errorf("validation error: %s MongoDB database not set", storage)
	}

	collection := readString(collectionKey)
	if collection == nil {
		return nil, fmt.Errorf("validation error: %s MongoDB collection not set", storage)
	}

	return &MongoDBDriverConfig{
		ConnectURL: *url,
		Database:   *db,
		Collection: *collection,
	}, nil
}

func readDynamoDBConfig(storage string) (*DynamoDBDriverConfig, error) {
	tableKeyName := fmt.Sprintf("%s_%s_DDB_TABLE_NAME", envStorageDriverPrefix, strings.ToUpper(storage))

	tableName := readString(tableKeyName)
	if tableName == nil {
		return nil, fmt.Errorf("validation error: %s DynamoDB table name not set", storage)
	}

	return &DynamoDBDriverConfig{
		TableName: *tableName,
	}, nil
}

func readS3Config(storage string) (*S3DriverConfig, error) {
	bucketKeyName := fmt.Sprintf("%s_%s_S3_BUCKET_NAME", envStorageDriverPrefix, strings.ToUpper(storage))
	bucketName := readString(bucketKeyName)
	if bucketName == nil {
		return nil, fmt.Errorf("validation error: %s s3 bucket name not set", storage)
	}

	privateKeyNameKey := fmt.Sprintf("%s_%s_S3_PRIVATE_KEY_NAME", envStorageDriverPrefix, strings.ToUpper(storage))
	privateKeyName := readString(privateKeyNameKey)
	if privateKeyName == nil {
		return nil, fmt.Errorf("validation error: %s s3 private key name not set", storage)
	}

	jwksNameKey := fmt.Sprintf("%s_%s_S3_JWKS_NAME", envStorageDriverPrefix, strings.ToUpper(storage))
	jwksName := readString(jwksNameKey)
	if jwksName == nil {
		return nil, fmt.Errorf("validation error: %s s3 jwks name not set", storage)
	}

	return &S3DriverConfig{
		BucketName:     *bucketName,
		PrivateKeyName: *privateKeyName,
		JWKSName:       *jwksName,
	}, nil
}

func readStorageDriverType(name string) (*string, error) {
	driverKey := fmt.Sprintf("%s_%s", envStorageDriverPrefix, strings.ToUpper(name))

	driver := readString(driverKey)
	if driver == nil {
		return nil, fmt.Errorf("validation error: storage driver %s not set", name)
	}

	// Check list of supported drivers
	if *driver != "dynamodb" && *driver != "mongodb" && *driver != "s3" {
		return nil, fmt.Errorf("validation error: storage driver %s not supported", *driver)
	}

	return driver, nil
}

func readString(name string) *string {
	val := os.Getenv(name)
	if len(val) == 0 {
		return nil
	}
	return &val
}

func readStringOrDefault(name string, defaultValue string) string {
	val := readString(name)
	if val == nil {
		return defaultValue
	}
	return *val
}

func readConfig() error {
	stage := readStringOrDefault(envStage, "dev")

	iss := readString(envIssClaim)
	if iss == nil {
		return fmt.Errorf("validation error: iss not set")
	}

	accessDurationString := readStringOrDefault(envAccessTokenLifetime, "60m")
	accessDuration, err := time.ParseDuration(accessDurationString)
	if err != nil {
		return fmt.Errorf("validation error: can not parse access token duration %s", accessDurationString)
	}

	refreshDurationString := readStringOrDefault(envRefreshTokenLifetime, "720h")
	refreshDuration, err := time.ParseDuration(refreshDurationString)
	if err != nil {
		return fmt.Errorf("validation error: can not parse refresh token duration %s", refreshDurationString)
	}

	keyLengthString := readString(envPrivateKeyLength)
	if keyLengthString == nil {
		return fmt.Errorf("validation error: key length not set")
	}
	keyLength, err := strconv.Atoi(*keyLengthString)
	if err != nil {
		return fmt.Errorf("validation error: key length %v", err)
	}

	usersStorageDriverType, err := readStorageDriverType("users")
	if err != nil {
		return err
	}

	keysStorageDriverType, err := readStorageDriverType("keys")
	if err != nil {
		return err
	}

	var usersStorageDriverCfg any
	switch *usersStorageDriverType {
	case "dynamodb":
		cfg, err := readDynamoDBConfig("users")
		if err != nil {
			return err
		}
		usersStorageDriverCfg = *cfg
	case "mongodb":
		cfg, err := readMongoDBConfig("users")
		if err != nil {
			return err
		}
		usersStorageDriverCfg = *cfg
	case "s3":
		cfg, err := readS3Config("users")
		if err != nil {
			return err
		}
		usersStorageDriverCfg = *cfg
	}

	var keysStorageDriverCfg any
	switch *keysStorageDriverType {
	case "dynamodb":
		cfg, err := readDynamoDBConfig("keys")
		if err != nil {
			return err
		}
		keysStorageDriverCfg = *cfg
	case "mongodb":
		cfg, err := readMongoDBConfig("keys")
		if err != nil {
			return err
		}
		keysStorageDriverCfg = *cfg
	case "s3":
		cfg, err := readS3Config("keys")
		if err != nil {
			return err
		}
		keysStorageDriverCfg = *cfg
	}

	Cfg = Config{
		Stage:                stage,
		TokenIssClaim:        *iss,
		TokenAccessDuration:  accessDuration,
		TokenRefreshDuration: refreshDuration,
		KeyLength:            keyLength,
		UserStorageDriver:    usersStorageDriverCfg,
		KeysStorageDriver:    keysStorageDriverCfg,
	}

	log.Println("Config", Cfg)

	return nil
}
