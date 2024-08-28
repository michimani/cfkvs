package output

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
)

type keyValueStore struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Status  string `json:"status"`
	ARN     string `json:"arn"`
}

type keyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type keyValueStoreSimple struct {
	ItemCount int32 `json:"itemCount"`
	TotalSize int64 `json:"totalSize"`
}

func toJson(data any) (any, error) {
	switch data := data.(type) {
	case *cfTypes.KeyValueStore:
		// CloudFront Key Value Store
		return keyValueStore{
			Id:      aws.ToString(data.Id),
			Name:    aws.ToString(data.Name),
			Comment: aws.ToString(data.Comment),
			Status:  aws.ToString(data.Status),
			ARN:     aws.ToString(data.ARN),
		}, nil

	case []cfTypes.KeyValueStore:
		// List of CloudFront Key Value Stores
		kvsList := make([]keyValueStore, 0, len(data))
		for _, kvs := range data {
			kvsList = append(kvsList, keyValueStore{
				Id:      aws.ToString(kvs.Id),
				Name:    aws.ToString(kvs.Name),
				Comment: aws.ToString(kvs.Comment),
				Status:  aws.ToString(kvs.Status),
				ARN:     aws.ToString(kvs.ARN),
			})
		}
		return kvsList, nil

	case []kvsTypes.ListKeysResponseListItem:
		// List of Items in the Key Value Store
		kvsItems := make([]keyValuePair, 0, len(data))
		for _, item := range data {
			kvsItems = append(kvsItems, keyValuePair{
				Key:   aws.ToString(item.Key),
				Value: aws.ToString(item.Value),
			})
		}
		return kvsItems, nil

	case *kvs.GetKeyOutput:
		// An item in the Key Value Store
		return keyValuePair{
			Key:   aws.ToString(data.Key),
			Value: aws.ToString(data.Value),
		}, nil

	case *kvs.PutKeyOutput:
		// Put an item in the Key Value Store
		return keyValueStoreSimple{
			ItemCount: aws.ToInt32(data.ItemCount),
			TotalSize: aws.ToInt64(data.TotalSizeInBytes),
		}, nil

	case *kvs.DeleteKeyOutput:
		// Delete an item in the Key Value Store
		return keyValueStoreSimple{
			ItemCount: aws.ToInt32(data.ItemCount),
			TotalSize: aws.ToInt64(data.TotalSizeInBytes),
		}, nil

	default:
		return nil, errors.New("unsupported data type")
	}
}

const jsonOutputIndent = "    " // 4 spaces

func RenderAsJson(data any) error {
	j, err := toJson(data)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(&j, "", jsonOutputIndent)
	if err != nil {
		return err
	}

	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n"))

	return nil
}
