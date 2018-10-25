package main

import (
	"fmt"
	"log"
	"os"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
)

// see struct_00.go for data structure

func ConvertAudio(ctx context.Context, sns events.SNSEvent) (string, error) {

	id := sns.Records[0].SNS.Message

	log.Printf("Processing id: %s\n", id)
	log.Printf("** REGION_AWS: %s \n", os.Getenv(REGION_AWS))
	log.Printf("** DB_TABLE_NAME: %s \n", os.Getenv(DB_TABLE_NAME))
	log.Printf("** SNS_TOPIC: %s, \n", os.Getenv(SNS_TOPIC))
	log.Printf("** S3_OUTPUT_BUCKET: %s \n", os.Getenv(S3_OUTPUT_BUCKET))
	log.Printf("** S3_OUTPUT_BUCKET_REGION: %s \n", os.Getenv(S3_OUTPUT_BUCKET_REGION))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(REGION_AWS)),
	})

	if (nil != err) {
		return "", err
		//log.Fatalf("%v\n", err)
	}

	ddbSess := dynamodb.New(sess)

	result, err := ddbSess.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv(DB_TABLE_NAME)),
		Key: map[string]* dynamodb.AttributeValue{
			"id": { S: aws.String(id), },
		},
	})
	if (nil != err) {
		return "", err
	}

	log.Printf("Dynamodb: GetItem: %v\n", result.Item)

	postItem := PostItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &postItem)
	if (nil != err) {
		//return "", err
		log.Fatalf("%v\n", err)
	}

	log.Printf("Dynamodb: UnmarshalMap: %v\n", postItem)

	if (len(postItem.Text) <= 0) {
		return "", errors.New("When you say nothing at all")
	}

	pollySess := polly.New(sess)
	pollyOutput, err := pollySess.SynthesizeSpeechWithContext(ctx, 
		&polly.SynthesizeSpeechInput{
			Text: aws.String(postItem.Text),
			VoiceId: aws.String(postItem.Voice),
			OutputFormat: aws.String("mp3"),
		},
	)
	if (nil != err) {
		return "", err
	}

	log.Printf("Polly: SynthesizeSpeech: %v\n", pollyOutput);

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(S3_OUTPUT_BUCKET_REGION)),
	})
	if (nil != err) {
		return "", err
	}

	audioFile := fmt.Sprintf("%s.mp3", postItem.Id)

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.UploadWithContext(ctx, 
		&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv(S3_OUTPUT_BUCKET)),
			Key: aws.String(audioFile),
			Body: pollyOutput.AudioStream,
		},
	)

	if (nil != err) {
		return "", err
	}

	log.Printf("S3: Upload: %s\n", audioFile)

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv(DB_TABLE_NAME)),
		Key: map[string]* dynamodb.AttributeValue{
			"id": {  S: aws.String(id) },
		},
		ExpressionAttributeValues: map[string]* dynamodb.AttributeValue{
			":stage": { S: aws.String("COMPLETED") },
		},
		UpdateExpression: aws.String("SET stage = :stage"),
		ReturnValues: aws.String("UPDATED_OLD"),
	}

	updateResult, err := ddbSess.UpdateItemWithContext(ctx, updateItemInput)
	if (nil != err) {
		return "", err
	}

	log.Printf("Dynamodb: UpdateItem: %v\n", updateResult)

	return audioFile, nil
}

func main() {
	lambda.Start(ConvertAudio)
}
