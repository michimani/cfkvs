package types

import (
	"errors"

	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
)

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (i *Item) Parse(o any) error {
	if i == nil {
		i = &Item{}
	}

	switch o := o.(type) {
	case *kvsTypes.ListKeysResponseListItem:
		i.Key = *o.Key
		i.Value = *o.Value
		return nil

	case *kvs.GetKeyOutput:
		i.Key = *o.Key
		i.Value = *o.Value
		return nil

	default:
		return errors.New("unexpected type")
	}
}

type ItemList []Item

func (il *ItemList) Parse(o *kvs.ListKeysOutput) error {
	if il == nil {
		il = &ItemList{}
	}

	for _, item := range o.Items {
		i := Item{}
		if err := i.Parse(item); err != nil {
			return err
		}
		*il = append(*il, i)
	}

	return nil
}
