package libs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	"github.com/michimani/cfkvs/types"
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
	DeleteKeyValueStore(ctx context.Context, params *cloudfront.DeleteKeyValueStoreInput, optFns ...func(*cloudfront.Options)) (*cloudfront.DeleteKeyValueStoreOutput, error)
	DescribeKeyValueStore(ctx context.Context, params *cloudfront.DescribeKeyValueStoreInput, optFns ...func(*cloudfront.Options)) (*cloudfront.DescribeKeyValueStoreOutput, error)
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
	Type() cfTypes.ImportSourceType
}

type KVSImportSourceS3 struct {
	Bucket string
	Key    string
}

func (s3 KVSImportSourceS3) ARN() string {
	return fmt.Sprintf("arn:aws:s3:::%s/%s", s3.Bucket, s3.Key)
}

func (s3 KVSImportSourceS3) Type() cfTypes.ImportSourceType {
	return cfTypes.ImportSourceTypeS3
}

func CreateKvs(ctx context.Context, c CloudFrontClient, name string, comment string, source KVSImportSource) (*cloudfront.CreateKeyValueStoreOutput, error) {
	input := &cloudfront.CreateKeyValueStoreInput{
		Name:    aws.String(name),
		Comment: aws.String(comment),
	}

	if source != nil {
		input.ImportSource = &cfTypes.ImportSource{
			SourceARN:  aws.String(source.ARN()),
			SourceType: source.Type(),
		}
	}

	return c.CreateKeyValueStore(ctx, input)
}

func DeleteKeyValueStore(ctx context.Context, c CloudFrontClient, kvsName string) error {
	eTag, err := getETagByCloudFront(ctx, c, kvsName)
	if err != nil {
		return err
	}

	input := &cloudfront.DeleteKeyValueStoreInput{
		Name:    aws.String(kvsName),
		IfMatch: eTag,
	}

	if out, err := c.DeleteKeyValueStore(ctx, input); err != nil {
		return err
	} else if out == nil {
		return fmt.Errorf("cloudfront.DeleteKeyValueStoreOutput is nil")
	}

	return nil
}

// DescribeKeyValueStore describes the key value store.
// The response of this function is a merge of CloudFront:DescribeKeyValueStore and CloudFrontKeyValueStore:DescribeKeyValueStore.
func DescribeKeyValueStore(ctx context.Context, cfc CloudFrontClient, kvsc CloudFrontKeyValueStoreClient, kvsName string) (*types.KeyValueStoreFull, error) {
	cIn := &cloudfront.DescribeKeyValueStoreInput{
		Name: aws.String(kvsName),
	}
	cOut, err := cfc.DescribeKeyValueStore(ctx, cIn)
	if err != nil {
		return nil, err
	}
	if cOut == nil {
		return nil, fmt.Errorf("cloudfront.DescribeKeyValueStoreOutput is nil")
	}
	if cOut.KeyValueStore == nil {
		return nil, fmt.Errorf("cloudfront.DescribeKeyValueStoreOutput.KeyValueStore is nil")
	}

	kvsIn := &kvs.DescribeKeyValueStoreInput{
		KvsARN: cOut.KeyValueStore.ARN,
	}
	kvsOut, err := kvsc.DescribeKeyValueStore(ctx, kvsIn)
	if err != nil {
		return nil, err
	}
	if kvsOut == nil {
		return nil, fmt.Errorf("cloudfrontkeyvaluestore.DescribeKeyValueStoreOutput is nil")
	}

	return &types.KeyValueStoreFull{
		ID:               aws.ToString(cOut.KeyValueStore.Id),
		ARN:              aws.ToString(cOut.KeyValueStore.ARN),
		Name:             aws.ToString(cOut.KeyValueStore.Name),
		Comment:          aws.ToString(cOut.KeyValueStore.Comment),
		Status:           aws.ToString(cOut.KeyValueStore.Status),
		ItemCount:        aws.ToInt32(kvsOut.ItemCount),
		TotalSizeInBytes: aws.ToInt64(kvsOut.TotalSizeInBytes),
		Created:          aws.ToTime(kvsOut.Created),
		LastModified:     aws.ToTime(kvsOut.LastModified),
		FailureReason:    aws.ToString(kvsOut.FailureReason),
		ETag:             aws.ToString(cOut.ETag),
	}, nil
}

func getETagByCloudFront(ctx context.Context, c CloudFrontClient, name string) (*string, error) {
	input := &cloudfront.DescribeKeyValueStoreInput{
		Name: aws.String(name),
	}
	out, err := c.DescribeKeyValueStore(ctx, input)
	if err != nil {
		return nil, err
	}

	return out.ETag, nil
}
