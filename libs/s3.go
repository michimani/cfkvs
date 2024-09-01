package libs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michimani/cfkvs/types"
)

func NewS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	c := s3.NewFromConfig(cfg)
	return c, nil
}

type S3Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func GetKeyValueStoreData(ctx context.Context, c S3Client, bucket, key string) (*types.KeyValueStoreData, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	out, err := c.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}

	kvsData := types.KeyValueStoreData{}
	if err := json.Unmarshal(bodyBytes, &kvsData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key value store data: %w", err)
	}

	return &kvsData, nil
}
