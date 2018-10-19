package main

import (
	"os"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	region := "us-east-2"
	if len(os.Args) > 1 {
		region = os.Args[1]
	}

	fmt.Printf("Buckets from %s region\n", region)

	sess, err := session.NewSession(&aws.Config{ Region: aws.String(region) })
	CheckError(err)

	s3Sess := s3.New(sess)
	output, err := s3Sess.ListBuckets(&s3.ListBucketsInput{})

	CheckError(err)

	for _, v := range output.Buckets {
		fmt.Printf("> %v\n", *v.Name)
	}
}
