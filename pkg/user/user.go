package user

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Credentials struct {
	Username string
	Password string
}

func GetCredentialsFromString(asString string) (*Credentials, error) {
	c := &Credentials{}
	err := json.Unmarshal([]byte(asString), c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func AddUser(creds Credentials) error {
	sess := session.Must(session.NewSession())
	ddbSvc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(creds.Username),
			},
			"password": {
				S: aws.String(creds.Password),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(username)"),
	}
	_, err := ddbSvc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(username string) (*Credentials, error) {
	return &Credentials{}, nil
}
