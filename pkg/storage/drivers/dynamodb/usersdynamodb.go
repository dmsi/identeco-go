package dynamodbdriver

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/storage"
)

type UsersDynamoDBDriver struct {
	lg    *slog.Logger
	ddb   *dynamodb.DynamoDB
	table string
}

func NewUsersStorage(lg *slog.Logger) (*UsersDynamoDBDriver, error) {
	cfg, ok := config.Cfg.UserStorageDriver.(config.DynamoDBDriverConfig)
	if !ok {
		return nil, fmt.Errorf("users driver configuration is not provided: dynamodb")
	}

	sess := session.New()

	return &UsersDynamoDBDriver{
		lg:    lg,
		ddb:   dynamodb.New(sess),
		table: cfg.TableName,
	}, nil
}

func (s UsersDynamoDBDriver) ReadUserData(username string) (*storage.UserData, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.table),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}

	item, err := s.ddb.GetItem(input)
	if err != nil {
		return nil, err
	}

	if item.Item == nil {
		return nil, fmt.Errorf("unable to read user: user %s not found", username)
	}

	user := &struct {
		Username string `dynamodbav:"username"`
		Hash     string `dynamodbav:"hash"`
		Token    string `dynamodbav:"refresh_token"`
	}{}

	err = dynamodbattribute.UnmarshalMap(item.Item, user)
	if err != nil {
		return nil, err
	}

	return &storage.UserData{
		Username: user.Username,
		Hash:     user.Hash,
	}, nil
}

func (s UsersDynamoDBDriver) WriteUserData(user storage.UserData) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.table),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"hash": {
				S: aws.String(user.Hash),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(username)"),
	}

	_, err := s.ddb.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
