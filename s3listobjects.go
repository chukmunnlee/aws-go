package main

import (
	"fmt"
	"log"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: bucket-name")
	}

	fmt.Printf("Bucket name: %s\n", os.Args[1])

	sess, err := session.NewSession(
			&aws.Config{ Region: aws.String("us-east-2") })
	CheckError(err)

	s3Sess := s3.New(sess)
	output, err := s3Sess.ListObjectsV2(
			&s3.ListObjectsV2Input{Bucket: aws.String(os.Args[1]) })
	CheckError(err)

	for _, v := range output.Contents {
		fmt.Printf("%v\n", *v.Key)
	}
}
