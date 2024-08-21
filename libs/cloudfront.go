package libs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func NewCloudFrontClient(ctx context.Context) (*cloudfront.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	c := cloudfront.NewFromConfig(cfg)
	return c, nil
}

type CloudFrontClient interface {
	ListKeyValueStores(ctx context.Context, params *cloudfront.ListKeyValueStoresInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListKeyValueStoresOutput, error)
	CreateKeyValueStore(ctx context.Context, params *cloudfront.CreateKeyValueStoreInput, optFns ...func(*cloudfront.Options)) (*cloudfront.CreateKeyValueStoreOutput, error)
}

func GetKeyValueStoreArn(ctx context.Context, c CloudFrontClient, kvsName string) (string, error) {
	input := &cloudfront.ListKeyValueStoresInput{}
	out, err := c.ListKeyValueStores(ctx, input)
	if err != nil {
		return "", err
	}

	for _, kvs := range out.KeyValueStoreList.Items {
		if *kvs.Name == kvsName {
			return *kvs.ARN, nil
		}
	}

	return "", fmt.Errorf("the key value store '%s' is not found", kvsName)
}

func ListKvs(ctx context.Context, c CloudFrontClient) ([]types.KeyValueStore, error) {
	input := &cloudfront.ListKeyValueStoresInput{}
	out, err := c.ListKeyValueStores(ctx, input)
	if err != nil {
		return nil, err
	}

	return out.KeyValueStoreList.Items, nil
}

func CreateKvs(ctx context.Context, c CloudFrontClient, name string) (*types.KeyValueStore, error) {
	input := &cloudfront.CreateKeyValueStoreInput{
		Name: &name,
	}
	out, err := c.CreateKeyValueStore(ctx, input)
	if err != nil {
		return nil, err
	}

	return out.KeyValueStore, nil
}
