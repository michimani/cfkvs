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

type KeyValueStoreData struct {
	Data []Item `json:"data"`
}

func (kd *KeyValueStoreData) ToItemList() *ItemList {
	il := &ItemList{
		Data:  []Item{},
		kvMap: map[string]*Item{},
	}

	for _, item := range kd.Data {
		il.Data = append(il.Data, item)
		il.kvMap[item.Key] = &item
	}

	return il
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

type ItemList struct {
	Data  []Item
	kvMap map[string]*Item
}

func (il *ItemList) Parse(o *kvs.ListKeysOutput) error {
	if il == nil {
		il = &ItemList{
			Data:  []Item{},
			kvMap: map[string]*Item{},
		}
	}

	if il.Data == nil {
		il.Data = []Item{}
	}

	if il.kvMap == nil {
		il.kvMap = map[string]*Item{}
	}

	for _, item := range o.Items {
		i := Item{}
		if err := i.Parse(&item); err != nil {
			return err
		}
		il.Data = append(il.Data, i)
		il.kvMap[i.Key] = &i
	}

	return nil
}

type ItemDiff struct {
	Before *Item
	After  *Item
}

type ItemListDiff struct {
	Add    []ItemDiff
	Update []ItemDiff
	Delete []ItemDiff
}

func (il *ItemList) Diff(after *ItemList, delete bool) *ItemListDiff {
	diff := &ItemListDiff{}

	// Update or Delete
	for _, before := range il.Data {
		after, ok := after.kvMap[before.Key]

		if !ok {
			if delete {
				// Delete
				diff.Delete = append(diff.Delete, ItemDiff{
					Before: &before,
					After:  nil,
				})
			}
			continue
		}

		// Update
		if before.Value != after.Value {
			diff.Update = append(diff.Update, ItemDiff{
				Before: &before,
				After:  after,
			})
		}
	}

	// Add
	for _, after := range after.Data {
		if _, ok := il.kvMap[after.Key]; !ok {
			diff.Add = append(diff.Add, ItemDiff{
				Before: nil,
				After:  &after,
			})
		}
	}

	return diff
}
