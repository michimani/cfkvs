package libs_test

import (
	"testing"

	"github.com/michimani/cfkvs/libs"
	"github.com/michimani/cfkvs/types"
	"github.com/stretchr/testify/assert"
)

func Test_GetKeyValueStoreDataFromFile(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		want    *types.KeyValueStoreData
		wantErr bool
	}{
		{
			name: "ok",
			path: "../testdata/valid.json",
			want: &types.KeyValueStoreData{
				Data: &[]types.Item{
					{
						Key:   "key-1",
						Value: "v 1",
					},
					{
						Key:   "key-2",
						Value: "value-2",
					},
					{
						Key:   "key-4",
						Value: "v 4",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "file not found",
			path:    "../testdata/notfound.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid json 1: invalid data structure",
			path:    "../testdata/invalid-1.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid json 2: key is empty",
			path:    "../testdata/invalid-2.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid json 3: value is empty",
			path:    "../testdata/invalid-3.json",
			want:    nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			got, err := libs.GetKeyValueStoreDataFromFile(c.path)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(got)
				return
			}

			asst.NoError(err)
			asst.Equal(*c.want, *got)
		})
	}
}
