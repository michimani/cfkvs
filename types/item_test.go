package types_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
	"github.com/michimani/cfkvs/types"
	"github.com/stretchr/testify/assert"
)

func Test_KeyValueStoreData_FromBytes(t *testing.T) {
	cases := []struct {
		name      string
		kvsd      *types.KeyValueStoreData
		input     []byte
		expect    types.KeyValueStoreData
		wantError bool
	}{
		{
			name:  "normal",
			kvsd:  &types.KeyValueStoreData{},
			input: []byte(`{"data":[{"key":"key1","value":"value1"},{"key":"key2","value":"value2"}]}`),
			expect: types.KeyValueStoreData{
				Data: &[]types.Item{
					{Key: "key1", Value: "value1"},
					{Key: "key2", Value: "value2"},
				},
			},
			wantError: false,
		},
		{
			name:  "empty",
			kvsd:  &types.KeyValueStoreData{},
			input: []byte(`{"data":[]}`),
			expect: types.KeyValueStoreData{
				Data: &[]types.Item{},
			},
			wantError: false,
		},
		{
			name:      "has empty key",
			kvsd:      &types.KeyValueStoreData{},
			input:     []byte(`{"data":[{"key":"","value":"value1"},{"key":"key2","value":"value2"}]}`),
			expect:    types.KeyValueStoreData{},
			wantError: true,
		},
		{
			name:      "has empty value",
			kvsd:      &types.KeyValueStoreData{},
			input:     []byte(`{"data":[{"key":"","value":""},{"key":"key2","value":"value2"}]}`),
			expect:    types.KeyValueStoreData{},
			wantError: true,
		},
		{
			name:      "invalid json: not json format",
			kvsd:      &types.KeyValueStoreData{},
			input:     []byte(`invalid`),
			expect:    types.KeyValueStoreData{},
			wantError: true,
		},
		{
			name:      "invalid json: invalid data structure 1",
			kvsd:      &types.KeyValueStoreData{},
			input:     []byte(`{"data":{}}`),
			expect:    types.KeyValueStoreData{},
			wantError: true,
		},
		{
			name:      "invalid json: invalid data structure 2",
			kvsd:      &types.KeyValueStoreData{},
			input:     []byte(`{"invalid":[]}`),
			expect:    types.KeyValueStoreData{},
			wantError: true,
		},
		{
			name:      "invalid json: invalid data structure 3",
			kvsd:      &types.KeyValueStoreData{},
			input:     []byte(`{"data":[{"k":"k1","v":"v1"}]}`),
			expect:    types.KeyValueStoreData{},
			wantError: true,
		},
		{
			name:      "nil KeyValueStoreData",
			kvsd:      nil,
			input:     []byte(`{"data":[{"key":"key1","value":"value1"},{"key":"key2","value":"value2"}]}`),
			expect:    types.KeyValueStoreData{},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.kvsd.FromBytes(c.input)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.Nil(err)
			asst.Equal(c.expect, *c.kvsd, c.kvsd)
		})
	}
}

func Test_KeyValueStoreData_ToItemList(t *testing.T) {
	cases := []struct {
		name       string
		kvsd       *types.KeyValueStoreData
		expectData []types.Item
	}{
		{
			name: "normal",
			kvsd: &types.KeyValueStoreData{
				Data: &[]types.Item{
					{Key: "key1", Value: "value1"},
					{Key: "key2", Value: "value2"},
				},
			},
			expectData: []types.Item{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"}},
		},
		{
			name: "empty",
			kvsd: &types.KeyValueStoreData{
				Data: &[]types.Item{},
			},
			expectData: []types.Item{},
		},
		{
			name:       "nil",
			kvsd:       nil,
			expectData: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			il := c.kvsd.ToItemList()

			if c.expectData == nil {
				asst.Nil(il)
				return
			}

			asst.NotNil(il)
			asst.Equal(c.expectData, il.Data)

			expectIl := types.NewItemList(c.expectData)
			asst.Equal(expectIl, il)
		})
	}
}

func Test_Item_Parse(t *testing.T) {
	cases := []struct {
		name      string
		item      *types.Item
		input     any
		expect    types.Item
		wantError bool
	}{
		{
			name:      "from kvsTypes.ListKeysResponseListItem",
			item:      &types.Item{},
			input:     &kvsTypes.ListKeysResponseListItem{Key: aws.String("key1"), Value: aws.String("value1")},
			expect:    types.Item{Key: "key1", Value: "value1"},
			wantError: false,
		},
		{
			name:      "from kvs.GetKeyOutput",
			item:      &types.Item{},
			input:     &kvs.GetKeyOutput{Key: aws.String("key2"), Value: aws.String("value2")},
			expect:    types.Item{Key: "key2", Value: "value2"},
			wantError: false,
		},
		{
			name:      "unexpected type",
			item:      &types.Item{},
			input:     "unexpected",
			expect:    types.Item{},
			wantError: true,
		},
		{
			name:      "nil item",
			item:      nil,
			input:     &kvsTypes.ListKeysResponseListItem{Key: aws.String("key1"), Value: aws.String("value1")},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.item.Parse(c.input)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.Nil(err)
			asst.Equal(c.expect.Key, c.item.Key)
			asst.Equal(c.expect.Value, c.item.Value)
		})
	}
}

func Test_NewItemList(t *testing.T) {
	cases := []struct {
		name       string
		items      []types.Item
		expectData []types.Item
	}{
		{
			name: "normal",
			items: []types.Item{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			expectData: []types.Item{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
		},
		{
			name:       "empty",
			items:      []types.Item{},
			expectData: []types.Item{},
		},
		{
			name:       "nil",
			items:      nil,
			expectData: []types.Item{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			il := types.NewItemList(c.items)

			if c.expectData == nil {
				asst.Nil(il)
				return
			}

			asst.NotNil(il)
			asst.Equal(c.expectData, il.Data)

			kvMap := types.GetKVMapFromItemList(il)
			for _, item := range c.expectData {
				asst.Equal(item, *kvMap[item.Key])
			}
		})
	}
}

func Test_ItemList_Parse(t *testing.T) {
	cases := []struct {
		name       string
		il         *types.ItemList
		input      *kvs.ListKeysOutput
		expectData []types.Item
		wantError  bool
	}{
		{
			name: "normal",
			il:   &types.ItemList{},
			input: &kvs.ListKeysOutput{
				Items: []kvsTypes.ListKeysResponseListItem{
					{Key: aws.String("key1"), Value: aws.String("value1")},
					{Key: aws.String("key2"), Value: aws.String("value2")},
				},
			},
			expectData: []types.Item{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			wantError: false,
		},
		{
			name: "nil ItemList",
			il:   nil,
			input: &kvs.ListKeysOutput{
				Items: []kvsTypes.ListKeysResponseListItem{
					{Key: aws.String("key1"), Value: aws.String("value1")},
					{Key: aws.String("key2"), Value: aws.String("value2")},
				},
			},
			expectData: nil,
			wantError:  true,
		},
		{
			name: "nil Items",
			il:   &types.ItemList{},
			input: &kvs.ListKeysOutput{
				Items: nil,
			},
			expectData: []types.Item{},
			wantError:  false,
		},
		{
			name:       "nil ListKeysOutput",
			il:         &types.ItemList{},
			input:      nil,
			expectData: []types.Item{},
			wantError:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.il.Parse(c.input)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.Nil(err)
			asst.Equal(c.expectData, c.il.Data)

			expectIl := types.NewItemList(c.expectData)
			asst.Equal(expectIl, c.il)
		})
	}
}

func Test_ItemList_Diff(t *testing.T) {
	ild1 := []types.Item{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
		{Key: "key3", Value: "value3"},
	}
	ild2 := []types.Item{
		{Key: "key1", Value: "v1"},
		{Key: "key2", Value: "value2"},
		{Key: "key4", Value: "value4"},
	}

	cases := []struct {
		name   string
		il1    *types.ItemList
		il2    *types.ItemList
		delete bool
		expect *types.ItemListDiff
	}{
		{
			name:   "no delete",
			il1:    types.NewItemList(ild1),
			il2:    types.NewItemList(ild2),
			delete: false,
			expect: &types.ItemListDiff{
				Add: []types.ItemDiff{
					{Before: nil, After: &types.Item{"key4", "value4"}},
				},
				Update: []types.ItemDiff{
					{Before: &types.Item{"key1", "value1"}, After: &types.Item{"key1", "v1"}},
				},
				Delete: []types.ItemDiff{},
			},
		},
		{
			name:   "delete",
			il1:    types.NewItemList(ild1),
			il2:    types.NewItemList(ild2),
			delete: true,
			expect: &types.ItemListDiff{
				Add: []types.ItemDiff{
					{Before: nil, After: &types.Item{"key4", "value4"}},
				},
				Update: []types.ItemDiff{
					{Before: &types.Item{"key1", "value1"}, After: &types.Item{"key1", "v1"}},
				},
				Delete: []types.ItemDiff{
					{Before: &types.Item{"key3", "value3"}, After: nil},
				},
			},
		},
		{
			name: "before is nil",
			il1:  nil,
			il2:  types.NewItemList(ild2),
			expect: &types.ItemListDiff{
				Add: []types.ItemDiff{
					{Before: nil, After: &types.Item{"key1", "v1"}},
					{Before: nil, After: &types.Item{"key2", "value2"}},
					{Before: nil, After: &types.Item{"key4", "value4"}},
				},
				Update: []types.ItemDiff{},
				Delete: []types.ItemDiff{},
			},
		},
		{
			name: "after is nil",
			il1:  types.NewItemList(ild1),
			il2:  nil,
			expect: &types.ItemListDiff{
				Add:    []types.ItemDiff{},
				Update: []types.ItemDiff{},
				Delete: []types.ItemDiff{},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			diff := c.il1.Diff(c.il2, c.delete)

			asst.NotNil(diff)
			asst.Equal(c.expect, diff)
		})
	}
}

func Test_ItemListDiff_PutList(t *testing.T) {
	cases := []struct {
		name   string
		ild    *types.ItemListDiff
		expect []types.Item
	}{
		{
			name: "normal",
			ild: &types.ItemListDiff{
				Add: []types.ItemDiff{
					{Before: nil, After: &types.Item{"key1", "value1"}},
					{Before: nil, After: &types.Item{"key2", "value2"}},
				},
				Update: []types.ItemDiff{
					{Before: &types.Item{"key3", "value3"}, After: &types.Item{"key3", "v3"}},
				},
				Delete: []types.ItemDiff{
					{Before: &types.Item{"key4", "value4"}, After: nil},
				},
			},
			expect: []types.Item{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "v3"},
			},
		},
		{
			name: "empty",
			ild: &types.ItemListDiff{
				Add:    []types.ItemDiff{},
				Update: []types.ItemDiff{},
				Delete: []types.ItemDiff{},
			},
			expect: []types.Item{},
		},
		{
			name:   "nil",
			ild:    nil,
			expect: []types.Item{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			putList := c.ild.PutList()

			asst.Equal(c.expect, putList)
		})
	}
}

func Test_ItemListDiff_DeleteList(t *testing.T) {
	cases := []struct {
		name   string
		ild    *types.ItemListDiff
		expect []types.Item
	}{
		{
			name: "normal",
			ild: &types.ItemListDiff{
				Add: []types.ItemDiff{
					{Before: nil, After: &types.Item{"key1", "value1"}},
					{Before: nil, After: &types.Item{"key2", "value2"}},
				},
				Update: []types.ItemDiff{
					{Before: &types.Item{"key3", "value3"}, After: &types.Item{"key3", "v3"}},
				},
				Delete: []types.ItemDiff{
					{Before: &types.Item{"key4", "value4"}, After: nil},
				},
			},
			expect: []types.Item{
				{Key: "key4", Value: "value4"},
			},
		},
		{
			name: "empty",
			ild: &types.ItemListDiff{
				Add:    []types.ItemDiff{},
				Update: []types.ItemDiff{},
				Delete: []types.ItemDiff{},
			},
			expect: []types.Item{},
		},
		{
			name:   "nil",
			ild:    nil,
			expect: []types.Item{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			deleteList := c.ild.DeleteList()

			asst.Equal(c.expect, deleteList)
		})
	}
}
