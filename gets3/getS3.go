package getS3

import (
    "fmt"
    "log"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct{}

func (s *S3Service) GetS3(bucket, key string) (string, error) {
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
    }))

    svc := s3.New(sess)

    // Get object metadata
    headInput := &s3.HeadObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    }
    result, err := svc.HeadObject(headInput)
    if err != nil {
        log.Println("Failed to get object metadata", err)
        return "", err
    }
    fmt.Println(*result.ContentType)

    if *result.ContentType != "video/mp4" {
        copyInput := &s3.CopyObjectInput{
            Bucket:     aws.String(bucket),
            Key:        aws.String(key),
            CopySource: aws.String(bucket + "/" + key),
            ContentType: aws.String("video/mp4"),
            MetadataDirective: aws.String("REPLACE"), // Replace metadata with the new one
        }
        _, err := svc.CopyObject(copyInput)
        if err != nil {
            log.Println("Failed to copy object", err)
            return "", err
        }
        fmt.Println("Content type changed to video/mp4")
    } else {
        fmt.Println("Content type is already video/mp4")
    }

    // Get presigned URL
    getObjectInput := &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    }
    req, _ := svc.GetObjectRequest(getObjectInput)
    urlStr, err := req.Presign(120 * time.Minute)
    if err != nil {
        log.Println("Failed to sign request", err)
        return "", err
    }

    log.Println("The URL is", urlStr)

    return urlStr, nil
}