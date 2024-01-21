package listS3

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "log"
)

type s3Service struct{}

func (s *s3Service) ListS3Objects(bucket, prefix string) ([]string, error) {
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
    }))

    svc := s3.New(sess)

    input := &s3.ListObjectsV2Input{
        Bucket: aws.String(bucket),
        Prefix: aws.String(prefix),
    }

    result, err := svc.ListObjectsV2(input)
    if err != nil {
        log.Println("Failed to list objects", err)
        return nil, err
    }

    keys := make([]string, len(result.Contents))
    for i, item := range result.Contents {
        keys[i] = *item.Key
    }

    return keys, nil
}