package output_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/types"
	"github.com/stretchr/testify/assert"
)

func Test_RenderAsTable(t *testing.T) {
	t1 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(1 * time.Hour)

	cases := []struct {
		name    string
		data    any
		expect  string
		wantErr bool
	}{
		{
			name: "ok: types.KVS",
			data: &types.KVS{Id: "id", Name: "name", Comment: "comment", Status: "status", ARN: "arn"},
			expect: `+----+------+---------+--------+-----+
| ID | NAME | COMMENT | STATUS | ARN |
+----+------+---------+--------+-----+
| id | name | comment | status | arn |
+----+------+---------+--------+-----+
`,
		},
		{
			name: "ok: types.KVSList",
			data: &types.KVSList{
				{Id: "id1", Name: "name1", Comment: "comment1", Status: "status1", ARN: "arn1"},
				{Id: "id2", Name: "name2", Comment: "comment2", Status: "status2", ARN: "arn2"},
			},
			expect: `+-----+-------+----------+---------+------+
| ID  | NAME  | COMMENT  | STATUS  | ARN  |
+-----+-------+----------+---------+------+
| id1 | name1 | comment1 | status1 | arn1 |
| id2 | name2 | comment2 | status2 | arn2 |
+-----+-------+----------+---------+------+
`,
		},
		{
			name: "ok: types.ItemList",
			data: &types.ItemList{
				Data: []types.Item{
					{Key: "key1", Value: "value1"},
					{Key: "key2", Value: "value2"},
				},
			},
			expect: `+------+--------+
| KEY  | VALUE  |
+------+--------+
| key1 | value1 |
| key2 | value2 |
+------+--------+
`,
		},
		{
			name: "ok: types.Item",
			data: &types.Item{Key: "key", Value: "value"},
			expect: `+-----+-------+
| KEY | VALUE |
+-----+-------+
| key | value |
+-----+-------+
`,
		},
		{
			name: "ok: types.KVSSimple",
			data: &types.KVSSimple{ItemCount: 1, TotalSize: 2},
			expect: `+-----------+------------------+
| ITEMCOUNT | TOTALSIZEINBYTES |
+-----------+------------------+
|         1 |                2 |
+-----------+------------------+
`,
		},
		{
			name: "ok: types.KeyValueStoreFull",
			data: &types.KeyValueStoreFull{
				ID:               "id",
				ARN:              "arn",
				Name:             "name",
				Comment:          "comment",
				Status:           "status",
				ItemCount:        1,
				TotalSizeInBytes: 2,
				Created:          t1,
				LastModified:     t2,
				FailureReason:    "failureReason",
				ETag:             "eTag",
			},
			expect: fmt.Sprintf(`+----+-----+------+---------+--------+-----------+------------------+-------------------------------+-------------------------------+---------------+------+
| ID | ARN | NAME | COMMENT | STATUS | ITEMCOUNT | TOTALSIZEINBYTES | CREATED                       | LASTMODIFIED                  | FAILUREREASON | ETAG |
+----+-----+------+---------+--------+-----------+------------------+-------------------------------+-------------------------------+---------------+------+
| id | arn | name | comment | status |         1 |                2 | %s | %s | failureReason | eTag |
+----+-----+------+---------+--------+-----------+------------------+-------------------------------+-------------------------------+---------------+------+
`, t1, t2),
		},
		{
			name: "ok: ItemListDiff",
			data: &types.ItemListDiff{
				Add: []types.ItemDiff{
					{Before: nil, After: &types.Item{Key: "key1", Value: "value1"}},
				},
				Update: []types.ItemDiff{
					{Before: &types.Item{Key: "key2", Value: "value2"}, After: &types.Item{Key: "key2", Value: "v2"}},
				},
				Delete: []types.ItemDiff{
					{Before: &types.Item{Key: "key3", Value: "value3"}, After: nil},
				},
			},
			expect: `
[ADDED] Following items will be added.
+---+------+--------+
| # | KEY  | VALUE  |
+---+------+--------+
| 1 | key1 | value1 |
+---+------+--------+

[UPDATED] Following items will be updated.
+---+------+--------------+-------------+
| # | KEY  | BEFORE VALUE | AFTER VALUE |
+---+------+--------------+-------------+
| 1 | key2 | value2       | v2          |
+---+------+--------------+-------------+

[DELETED] Following items will be deleted.
+---+------+--------+
| # | KEY  | VALUE  |
+---+------+--------+
| 1 | key3 | value3 |
+---+------+--------+
`,
		},
		{
			name: "ok: ItemListDiff (empty)",
			data: &types.ItemListDiff{
				Add:    []types.ItemDiff{},
				Update: []types.ItemDiff{},
				Delete: []types.ItemDiff{},
			},
			expect: `
[ADDED] No items will be added.

[UPDATED] No items will be updated.

[DELETED] No items will be deleted.
`,
		},
		{
			name:    "invalid data",
			data:    make(chan int),
			wantErr: true,
		},
	}

	out := new(bytes.Buffer)

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			tt.Cleanup(func() {
				out.Truncate(0)
			})

			asst := assert.New(tt)

			err := output.RenderAsTable(c.data, out)
			if c.wantErr {
				asst.Error(err)
				asst.Empty(out.String())
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, out.String())
		})
	}
}
