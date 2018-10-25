package main

import (
	//"fmt"
	"log"
	"os"
	"strings"
	//"encoding/json"
	"context"
	//"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/aws/aws-lambda-go/events"
)

const STAR string = "*"

func mkPostItem(v map[string] *dynamodb.AttributeValue) PostItem {
	return PostItem{
		Id: *v["id"].S,
		Stage: *v["stage"].S,
		Text: *v["text"].S,
		Voice: *v["voice"].S,
	}
}

func GetPosts(ctx context.Context, event GetPostEvent) ([]PostItem, error) {
//func GetPosts() {
	id := event.PostId

	log.Printf("Getting id: %s\n", id)
	log.Printf("** REGION_AWS: %s \n", os.Getenv(REGION_AWS))
	log.Printf("** DB_TABLE_NAME: %s \n", os.Getenv(DB_TABLE_NAME))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(REGION_AWS)),
	})
	if nil != err {
		return nil, err
		//log.Fatalf("Error: %v\n", err)
	}

	ddbSess := dynamodb.New(sess)

	result := []PostItem{}

	if 0 == strings.Compare(id, STAR) {
		output, err := ddbSess.ScanWithContext(ctx, 
			&dynamodb.ScanInput{
				TableName: aws.String(os.Getenv(DB_TABLE_NAME)),
			},
		)

		if nil != err {
			return nil, err
			//log.Fatalf("Error: %v\n", err)
		}
		log.Printf("Dynamodb: Scan: %v\n", len(output.Items))

		for _, v := range output.Items {
			result = append(result, mkPostItem(v))
		}
	} else {
		output, err := ddbSess.GetItemWithContext(ctx, 
			&dynamodb.GetItemInput{
				TableName: aws.String(os.Getenv(DB_TABLE_NAME)),
				Key: map[string] *dynamodb.AttributeValue{
					"id": { S: aws.String(id) },
				},
			},
		)
		if (nil != err) {
			return nil, err
			//log.Fatalf("Error: %v\n", err)
		}

		log.Printf("Dynamodb: GetItem: %v\n", output.Item)
		if len(output.Item) > 0 {
			result = append(result, mkPostItem(output.Item))
		}
	}

	/*
	j, err := json.Marshal(result)
	if nil != err {
		return "", err
		//log.Fatalf("Error: %v\n", err)
	}

	log.Printf("json: Marshal: %v\n", string(j))
	*/

	return result, nil
}

func main() {
	lambda.Start(GetPosts)
	//GetPosts()
}
