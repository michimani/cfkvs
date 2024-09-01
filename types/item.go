package types

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
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
		i.Key = aws.ToString(o.Key)
		i.Value = aws.ToString(o.Value)
		return nil

	case *kvs.GetKeyOutput:
		i.Key = aws.ToString(o.Key)
		i.Value = aws.ToString(o.Value)
		return nil

	default:
		return fmt.Errorf("failed to parse Item due to unexpected type: %s", reflect.TypeOf(o).String())
	}
}

type ItemList []Item

func (il *ItemList) Parse(o *kvs.ListKeysOutput) error {
	if il == nil {
		il = &ItemList{}
	}

	for _, item := range o.Items {
		i := Item{}
		if err := i.Parse(&item); err != nil {
			return err
		}
		*il = append(*il, i)
	}

	return nil
}
