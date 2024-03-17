package usersdynamodb

import (
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/storage"
)

type UsersStorage struct {
	lg    *slog.Logger
	ddb   *dynamodb.DynamoDB
	table string
}

func New(lg *slog.Logger, table string) *UsersStorage {
	sess := session.New()

	return &UsersStorage{
		lg:    lg,
		ddb:   dynamodb.New(sess),
		table: table,
	}
}

func wrap(name string, err error) error {
	return e.Wrap("storage.dynamodb.userdynamodb."+name, err)
}

func (u *UsersStorage) ReadUserData(username string) (*storage.UserData, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(u.table),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}

	item, err := u.ddb.GetItem(input)
	if err != nil {
		return nil, wrap("ReadUserData", err)
	}

	if item.Item == nil {
		return nil, wrap("ReadUserData", errors.New("user not found"))
	}

	user := &struct {
		Username string `dynamodbav:"username"`
		Hash     string `dynamodbav:"hash"`
		Token    string `dynamodbav:"refresh_token"`
	}{}

	err = dynamodbattribute.UnmarshalMap(item.Item, user)
	if err != nil {
		return nil, wrap("ReadUserData", err)
	}

	return &storage.UserData{
		Username: user.Username,
		Hash:     user.Hash,
	}, nil
}

func (u *UsersStorage) WriteUserData(user storage.UserData) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(u.table),
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

	_, err := u.ddb.PutItem(input)
	if err != nil {
		return wrap("WriteUserData", err)
	}

	return nil
}
