package main

import (
	//"fmt"
	"os"
	"log"
	//"io/ioutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	if (len(os.Args) < 3) {
		log.Fatalln("usage: bucket-name file-name");
	}

	bucketName := os.Args[1]
	fileName := os.Args[2]

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})
	CheckError(err)

	reader, err := os.Open(fileName)
	CheckError(err)
	defer reader.Close()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileName),
		Body: reader,
	})
	CheckError(err)

	log.Printf("Uploaded file %s to bucket %s\n", fileName, bucketName)
}
