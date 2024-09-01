package libs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func ListKvs(ctx context.Context, c CloudFrontClient) (*cloudfront.ListKeyValueStoresOutput, error) {
	input := &cloudfront.ListKeyValueStoresInput{}
	return c.ListKeyValueStores(ctx, input)
}

type KVSImportSource interface {
	ARN() string
	Type() types.ImportSourceType
}

type KVSImportSourceS3 struct {
	Bucket string
	Key    string
}

func (s3 KVSImportSourceS3) ARN() string {
	return fmt.Sprintf("arn:aws:s3:::%s/%s", s3.Bucket, s3.Key)
}

func (s3 KVSImportSourceS3) Type() types.ImportSourceType {
	return types.ImportSourceTypeS3
}

func CreateKvs(ctx context.Context, c CloudFrontClient, name string, comment string, source KVSImportSource) (*cloudfront.CreateKeyValueStoreOutput, error) {
	input := &cloudfront.CreateKeyValueStoreInput{
		Name:    aws.String(name),
		Comment: aws.String(comment),
	}

	if source != nil {
		input.ImportSource = &types.ImportSource{
			SourceARN:  aws.String(source.ARN()),
			SourceType: source.Type(),
		}
	}

	return c.CreateKeyValueStore(ctx, input)
}
