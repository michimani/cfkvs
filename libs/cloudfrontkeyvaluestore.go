package libs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	"github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
)

func NewCloudFrontKeyValueStoreClient(ctx context.Context) (*kvs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	c := kvs.NewFromConfig(cfg)
	return c, nil
}

type CloudFrontKeyValueStoreClient interface {
	ListKeys(ctx context.Context, params *kvs.ListKeysInput, optFns ...func(*kvs.Options)) (*kvs.ListKeysOutput, error)
	GetKey(ctx context.Context, params *kvs.GetKeyInput, optFns ...func(*kvs.Options)) (*kvs.GetKeyOutput, error)
	PutKey(ctx context.Context, params *kvs.PutKeyInput, optFns ...func(*kvs.Options)) (*kvs.PutKeyOutput, error)
	DeleteKey(ctx context.Context, params *kvs.DeleteKeyInput, optFns ...func(*kvs.Options)) (*kvs.DeleteKeyOutput, error)
	DescribeKeyValueStore(ctx context.Context, params *kvs.DescribeKeyValueStoreInput, optFns ...func(*kvs.Options)) (*kvs.DescribeKeyValueStoreOutput, error)
}

func ListItems(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN string) ([]types.ListKeysResponseListItem, error) {
	input := &kvs.ListKeysInput{
		KvsARN: aws.String(kvsARN),
	}
	out, err := c.ListKeys(ctx, input)
	if err != nil {
		return nil, err
	}

	return out.Items, nil
}

func GetItem(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN, key string) (*kvs.GetKeyOutput, error) {
	input := &kvs.GetKeyInput{
		KvsARN: aws.String(kvsARN),
		Key:    aws.String(key),
	}
	out, err := c.GetKey(ctx, input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func PutItem(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN, key, value string) (*kvs.PutKeyOutput, error) {
	eTag, err := getETag(ctx, c, kvsARN)
	if err != nil {
		return nil, err
	}

	input := &kvs.PutKeyInput{
		IfMatch: eTag,
		KvsARN:  aws.String(kvsARN),
		Key:     aws.String(key),
		Value:   aws.String(value),
	}
	out, err := c.PutKey(ctx, input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func DeleteItem(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN, key string) (*kvs.DeleteKeyOutput, error) {
	eTag, err := getETag(ctx, c, kvsARN)
	if err != nil {
		return nil, err
	}

	input := &kvs.DeleteKeyInput{
		IfMatch: eTag,
		KvsARN:  aws.String(kvsARN),
		Key:     aws.String(key),
	}

	out, err := c.DeleteKey(ctx, input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func getETag(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN string) (*string, error) {
	input := &kvs.DescribeKeyValueStoreInput{
		KvsARN: aws.String(kvsARN),
	}
	out, err := c.DescribeKeyValueStore(ctx, input)
	if err != nil {
		return nil, err
	}

	return out.ETag, nil
}
