package libs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
	"github.com/michimani/cfkvs/types"
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
	UpdateKeys(ctx context.Context, params *kvs.UpdateKeysInput, optFns ...func(*kvs.Options)) (*kvs.UpdateKeysOutput, error)
	DescribeKeyValueStore(ctx context.Context, params *kvs.DescribeKeyValueStoreInput, optFns ...func(*kvs.Options)) (*kvs.DescribeKeyValueStoreOutput, error)
}

func ListItems(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN string) (*types.ItemList, error) {
	input := &kvs.ListKeysInput{
		KvsARN: aws.String(kvsARN),
	}

	items := []types.Item{}
	for {
		out, err := c.ListKeys(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, item := range out.Items {
			items = append(items, types.Item{
				Key:   aws.ToString(item.Key),
				Value: aws.ToString(item.Value),
			})
		}

		input.NextToken = out.NextToken
		if input.NextToken == nil {
			break
		}
	}

	itemList := types.NewItemList(items)

	return itemList, nil
}

func GetItem(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN, key string) (*kvs.GetKeyOutput, error) {
	input := &kvs.GetKeyInput{
		KvsARN: aws.String(kvsARN),
		Key:    aws.String(key),
	}
	return c.GetKey(ctx, input)
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
	return c.PutKey(ctx, input)
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

	return c.DeleteKey(ctx, input)
}

func SyncItems(ctx context.Context, c CloudFrontKeyValueStoreClient, kvsARN string, putList, deleteList []types.Item) (*kvs.UpdateKeysOutput, error) {
	puts := []kvsTypes.PutKeyRequestListItem{}
	for _, item := range putList {
		puts = append(puts, kvsTypes.PutKeyRequestListItem{
			Key:   aws.String(item.Key),
			Value: aws.String(item.Value),
		})
	}

	delete := []kvsTypes.DeleteKeyRequestListItem{}
	for _, item := range deleteList {
		delete = append(delete, kvsTypes.DeleteKeyRequestListItem{
			Key: aws.String(item.Key),
		})
	}

	eTag, err := getETag(ctx, c, kvsARN)
	if err != nil {
		return nil, err
	}

	input := &kvs.UpdateKeysInput{
		KvsARN:  aws.String(kvsARN),
		Puts:    puts,
		Deletes: delete,
		IfMatch: eTag,
	}

	return c.UpdateKeys(ctx, input)
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
