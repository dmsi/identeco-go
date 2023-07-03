package userdata

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dmsi/identeco/pkg/lib/e"
	"github.com/dmsi/identeco/pkg/storage"
)

// TODO: logger interface
type UserDataStorage struct {
	ddb   *dynamodb.DynamoDB
	table string
}

func New(table string) *UserDataStorage {
	sess := session.New()

	return &UserDataStorage{
		ddb:   dynamodb.New(sess),
		table: table,
	}
}

func op(name string) string {
	return "storage.dynamodb.userdata." + name
}

func (u *UserDataStorage) ReadUserData(username string) (*storage.UserData, error) {
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
		return nil, err
	}

	// type userInfo struct {
	// 	Username string `dynamodbav:"username"`
	// 	Hash     string `dynamodbav:"hash"`
	// 	Token    string `dynamodbav:"refresh_token"`
	// }

	user := &struct {
		Username string `dynamodbav:"username"`
		Hash     string `dynamodbav:"hash"`
		Token    string `dynamodbav:"refresh_token"`
	}{}

	err = dynamodbattribute.UnmarshalMap(item.Item, user)
	if err != nil {
		return nil, e.Wrap(op("ReadUserData"), err)
	}

	return &storage.UserData{
		Username:     user.Username,
		Hash:         user.Hash,
		RefreshToken: user.Token,
	}, nil

	// fmt.Printf("Item >>> %v\n", item.Item)
	// uu := &userInfo{}
	// err = dynamodbattribute.UnmarshalMap(item.Item, uu)
	// fmt.Printf("uu %v, err: %v\n", *uu, err)

	// user := &storage.UserData{}
	// err = dynamodbattribute.UnmarshalMap(item.Item, user)
	// if err != nil {
	// 	return nil, err
	// }

	// return user, nil
}

func (u *UserDataStorage) WriteUserData(username string, user storage.UserData) error {
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
		return e.Wrap(op("WriteUserData"), err)
	}

	return nil
}
