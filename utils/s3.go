package utils

import (
	"context"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func PutObject(s3Client *s3.Client) {

}

func getObject(s3Client *s3.Client) {

}

func CopyObject(s3Client *s3.Client) {

}

func DeleteObject(s3Client *s3.Client, bucket string, key string) error {
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

func PresignPutObject(s3Client *s3.Client, bucket string, key string, contentType string) (string, error) {
	// create presign client to get presign URL for download
	presignClient := s3.NewPresignClient(s3Client)

	s3PresignData, err := presignClient.PresignPutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
		},
		s3.WithPresignExpires(15*time.Minute), // URL expires in 15 mins
	)

	if err != nil {
		return "", err
	}

	return s3PresignData.URL, nil
}

func PresignGetObject(s3Client *s3.Client, bucket string, key string) (string, error) {

	// create presign client to get presign URL for download
	presignClient := s3.NewPresignClient(s3Client)

	s3PresignData, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("attachment; filename=\"" + path.Base(key) + "\""),
	})

	if err != nil {
		return "", err
	}

	return s3PresignData.URL, nil

}
