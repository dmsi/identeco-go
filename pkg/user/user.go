package user

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string
	Password string
	Hash     string
}

func GetCredentialsFromString(asString string) (*Credentials, error) {
	c := &Credentials{}
	err := json.Unmarshal([]byte(asString), c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func hashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return aws.String(string(hash)), nil
}

func comparePassword(password, hash string) (*bool, error) {
	fmt.Printf(">>> comparing %v vs %v\n", hash, password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return nil, err
	}

	return aws.Bool(true), nil
}

func AddUser(creds Credentials) error {
	sess := session.Must(session.NewSession())
	ddbSvc := dynamodb.New(sess)

	hash, err := hashPassword(creds.Password)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(creds.Username),
			},
			// "password": {
			// 	S: aws.String(creds.Password),
			// },
			"hash": {
				S: hash,
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(username)"),
	}
	_, err = ddbSvc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(username string) (*Credentials, error) {
	// return &Credentials{}, nil
	sess := session.Must(session.NewSession())
	ddbSvc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}

	item, err := ddbSvc.GetItem(input)
	if err != nil {
		return nil, err
	}

	// item.Item
	creds := &Credentials{}
	// dynamodbattribute.UnmarhalMap(item.Item, creds)
	err = dynamodbattribute.UnmarshalMap(item.Item, creds)
	if err != nil {
		return nil, err
	}

	fmt.Printf(">>>> !!!! ----> %v\n", creds)

	return &Credentials{}, nil
}

// Returns error in case when password does not match
// Or when user not found
// Or when an error happened? TODO bool + error?
func VerifyPassword(creds Credentials) (*bool, error) {
	sess := session.Must(session.NewSession())
	ddbSvc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(creds.Username),
			},
		},
	}

	output, err := ddbSvc.GetItem(input)
	if err != nil {
		return nil, err
	}

	// User not found
	if output.Item == nil {
		return aws.Bool(false), nil
	}

	result := &Credentials{}
	err = dynamodbattribute.UnmarshalMap(output.Item, result)
	if err != nil {
		return nil, err
	}

	match, err := comparePassword(creds.Password, result.Hash)
	if err != nil || !*match {
		return aws.Bool(false), nil
	}

	return aws.Bool(true), nil
}

func VerifyUser(username string) (*bool, error) {
	_, err := GetUser(username)
	if err != nil {
		return nil, err
	}
	return aws.Bool(true), nil
}
