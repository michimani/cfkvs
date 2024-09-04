package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
)

type KVS struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Status  string `json:"status"`
	ARN     string `json:"arn"`
}

type KeyValueStoreFull struct {
	ID               string    `json:"id"`
	ARN              string    `json:"arn"`
	Name             string    `json:"name"`
	Comment          string    `json:"comment"`
	Status           string    `json:"status"`
	ItemCount        int32     `json:"itemCount"`
	TotalSizeInBytes int64     `json:"totalSizeInBytes"`
	Created          time.Time `json:"created"`
	LastModified     time.Time `json:"lastModified"`
	FailureReason    string    `json:"failureReason"`
	ETag             string    `json:"eTag"`
}

type KVSSimple struct {
	ItemCount int32 `json:"itemCount"`
	TotalSize int64 `json:"totalSize"`
}

func (k *KVS) Parse(o any) error {
	if k == nil {
		k = &KVS{}
	}

	switch o := o.(type) {
	case *cf.CreateKeyValueStoreOutput:
		k.Id = aws.ToString(o.KeyValueStore.Id)
		k.Name = aws.ToString(o.KeyValueStore.Name)
		k.Comment = aws.ToString(o.KeyValueStore.Comment)
		k.Status = aws.ToString(o.KeyValueStore.Status)
		k.ARN = aws.ToString(o.KeyValueStore.ARN)
		return nil

	case *cf.DescribeKeyValueStoreOutput:
		k.Id = aws.ToString(o.KeyValueStore.Id)
		k.Name = aws.ToString(o.KeyValueStore.Name)
		k.Comment = aws.ToString(o.KeyValueStore.Comment)
		k.Status = aws.ToString(o.KeyValueStore.Status)
		k.ARN = aws.ToString(o.KeyValueStore.ARN)
		return nil

	case *cfTypes.KeyValueStore:
		k.Id = aws.ToString(o.Id)
		k.Name = aws.ToString(o.Name)
		k.Comment = aws.ToString(o.Comment)
		k.Status = aws.ToString(o.Status)
		k.ARN = aws.ToString(o.ARN)
		return nil

	default:
		return fmt.Errorf("failed to parse KVS due to unexpected type: %s", reflect.TypeOf(o).String())
	}
}

type KVSList []KVS

func (kl *KVSList) Parse(o *cf.ListKeyValueStoresOutput) error {
	if kl == nil {
		kl = &KVSList{}
	}

	for _, kvs := range o.KeyValueStoreList.Items {
		k := KVS{}
		if err := k.Parse(&kvs); err != nil {
			return err
		}
		*kl = append(*kl, k)
	}

	return nil
}

func (ks *KVSSimple) Parse(o any) error {
	if ks == nil {
		ks = &KVSSimple{}
	}

	switch o := o.(type) {
	case *kvs.PutKeyOutput:
		ks.ItemCount = aws.ToInt32(o.ItemCount)
		ks.TotalSize = aws.ToInt64(o.TotalSizeInBytes)
		return nil

	case *kvs.DeleteKeyOutput:
		ks.ItemCount = aws.ToInt32(o.ItemCount)
		ks.TotalSize = aws.ToInt64(o.TotalSizeInBytes)
		return nil

	case *kvs.UpdateKeysOutput:
		ks.ItemCount = aws.ToInt32(o.ItemCount)
		ks.TotalSize = aws.ToInt64(o.TotalSizeInBytes)
		return nil

	default:
		return fmt.Errorf("failed to parse KVSSimple due to unexpected type: %s", reflect.TypeOf(o).String())
	}
}
