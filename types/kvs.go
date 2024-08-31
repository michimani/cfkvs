package types

import (
	"errors"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
)

type KVS struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Status  string `json:"status"`
	ARN     string `json:"arn"`
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
		k.Id = *o.KeyValueStore.Id
		k.Name = *o.KeyValueStore.Name
		k.Comment = *o.KeyValueStore.Comment
		k.Status = *o.KeyValueStore.Status
		k.ARN = *o.KeyValueStore.ARN
		return nil

	case *cf.DescribeKeyValueStoreOutput:
		k.Id = *o.KeyValueStore.Id
		k.Name = *o.KeyValueStore.Name
		k.Comment = *o.KeyValueStore.Comment
		k.Status = *o.KeyValueStore.Status
		k.ARN = *o.KeyValueStore.ARN
		return nil

	default:
		return errors.New("unexpected type")
	}
}

type KVSList []KVS

func (kl *KVSList) Parse(o *cf.ListKeyValueStoresOutput) error {
	if kl == nil {
		kl = &KVSList{}
	}

	for _, kvs := range o.KeyValueStoreList.Items {
		k := KVS{}
		if err := k.Parse(kvs); err != nil {
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
		ks.ItemCount = *o.ItemCount
		ks.TotalSize = *o.TotalSizeInBytes
		return nil

	case *kvs.DeleteKeyOutput:
		ks.ItemCount = *o.ItemCount
		ks.TotalSize = *o.TotalSizeInBytes
		return nil

	default:
		return errors.New("unexpected type")
	}
}
