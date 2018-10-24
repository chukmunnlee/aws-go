package main

import (
	//"fmt"
	"log"
	"os"
	"context"
	"github.com/satori/go.uuid"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-lambda-go/lambda"
)

// see struct_00.go for data structure

func NewPost(ctx context.Context, newPostEvent NewPostEvent) (string, error) {
	uuid := uuid.Must(uuid.NewV1())

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(REGION_AWS)),
	})
	if (nil != err) {
		return "", err
	}

	log.Printf("Generating id: %s\n", uuid.String())
	log.Printf("** REGION_AWS: %s, \n", os.Getenv(REGION_AWS))
	log.Printf("** DB_TABLE_NAME: %s, \n", os.Getenv(DB_TABLE_NAME))
	log.Printf("** SNS_TOPIC: %s, \n", os.Getenv(SNS_TOPIC))

	postItem := PostItem{
		Id: uuid.String(),
		Status: "PROCESSING",
		Text: newPostEvent.Text,
		Voice: newPostEvent.Voice,
	}

	av, err := dynamodbattribute.MarshalMap(postItem)
	if (nil != err) {
		return "", err
	}

	ddbSess := dynamodb.New(sess)

	putItem := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(os.Getenv(DB_TABLE_NAME)),
	}

	_, err = ddbSess.PutItemWithContext(ctx, putItem)
	if (nil != err) {
		return "", err
	}

	log.Printf("DynamoDB: PutItem %s \n", uuid)

	snsSess := sns.New(sess)
	_, err = snsSess.PublishWithContext(ctx, &sns.PublishInput{
		Message: aws.String(uuid.String()),
		TopicArn: aws.String(os.Getenv(SNS_TOPIC)),
	})
	if (nil != err) {
		return "", err
	}

	log.Printf("SNS: Publish: %s\n" , uuid)

	return uuid.String(), nil
}

func main() {
	lambda.Start(NewPost)
}
