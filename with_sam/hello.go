package main

import (
	"fmt"
	"log"
	"time"
	"strings"
	"errors"
	"github.com/satori/go.uuid"
	"os"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	HELLO_AWS_REGION = "HELLO_AWS_REGION"
	HELLO_BUCKET_NAME = "HELLO_BUCKET_NAME"
)

type HelloRequest struct {
	UserId string `json:"userId"`
	Message string `json:"message"`
}

type HelloResponse struct {
	ResponseId string `json:"respondId"`
	Time time.Time `json:"time"`
	Message string `json:"message"`
	Status int `json:"statusCode"`
}

func CheckError(err error) {
	if nil != err {
		log.Fatalf("Error: %v\n", err)
	}
}

func HandleRequest(ctx context.Context, req HelloRequest) (HelloResponse, error) {

	var resp  HelloResponse

	if ((len(req.Message) <= 0) || (len(req.UserId) <= 0)) {
		return resp, errors.New("Invalid request")
	}

	respId := uuid.Must(uuid.NewV4())
	prefix := respId.String()[:6]
	currTime := time.Now()
	fileName := fmt.Sprintf("%s-%s", prefix, req.UserId)

	log.Printf("AWS Region: %s\n", os.Getenv(HELLO_AWS_REGION))
	log.Printf("Bucket name: %s\n", os.Getenv(HELLO_BUCKET_NAME))
	log.Printf("Response id: %s\n", respId)
	log.Printf("File name: %s\n", fileName)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(HELLO_AWS_REGION)),
	})
	CheckError(err)

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(os.Getenv(HELLO_BUCKET_NAME)),
		Key: aws.String(fileName),
		Body: strings.NewReader(fmt.Sprintf("Date: %v: Message: %s", currTime, req.Message)),
	})

	resp = HelloResponse{ 
		ResponseId: respId.String(), 
		Time: currTime,
		Message: fmt.Sprintf("Hello %s: Filename: %s", req.UserId, fileName), 
		Status: 200,
	}

	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
	/*
	resp, err := HandleRequest(HelloRequest{
		UserId: "fred",
		Message: "This is a message",
	})

	if nil != err {
		log.Fatalf("Error: %s\n", err)
	}

	log.Printf("RESULT: %v\n", resp)
	*/
}
