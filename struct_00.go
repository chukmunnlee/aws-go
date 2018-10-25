package main

const (
	REGION_AWS string = "REGION_AWS"
	SNS_TOPIC string = "SNS_TOPIC"
	DB_TABLE_NAME string = "DB_TABLE_NAME"
	TEST_ID string = "7e71d533-d751-11e8-99da-7ee0c634495c"
	S3_OUTPUT_BUCKET string = "S3_OUTPUT_BUCKET"
	S3_OUTPUT_BUCKET_REGION string = "S3_OUTPUT_BUCKET_REGION"
)

type PostItem struct {
	Id string `json:"id"`
	Stage string `json:"stage"`
	Text string `json:"text"`
	Voice string `json:"voice"`
}

type NewPostEvent struct {
	Text string `json:"text"`
	Voice string `json:"voice"`
}

type GetPostEvent struct {
	PostId string `json:"postId"`
}
