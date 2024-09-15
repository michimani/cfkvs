package libs_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	"github.com/michimani/cfkvs/libs"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_getETag(t *testing.T) {
	cases := []struct {
		name    string
		kvscOut struct {
			Out   *kvs.DescribeKeyValueStoreOutput
			Error error
		}
		kvsARN  string
		expect  *string
		wantErr bool
	}{
		{
			name: "ok",
			kvscOut: struct {
				Out   *kvs.DescribeKeyValueStoreOutput
				Error error
			}{
				Out: &kvs.DescribeKeyValueStoreOutput{
					ETag: aws.String("dummy_etag"),
				},
				Error: nil,
			},
			kvsARN:  "dummy_arn",
			expect:  aws.String("dummy_etag"),
			wantErr: false,
		},
		{
			name: "failed to describe key value store",
			kvscOut: struct {
				Out   *kvs.DescribeKeyValueStoreOutput
				Error error
			}{
				Out:   nil,
				Error: assert.AnError,
			},
			kvsARN:  "dummy_arn",
			expect:  nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			ctrl := gomock.NewController(tt)

			m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
			m.EXPECT().
				DescribeKeyValueStore(gomock.Any(), gomock.Any()).
				Return(c.kvscOut.Out, c.kvscOut.Error)

			got, err := libs.Exported_getETag(context.Background(), m, c.kvsARN)
			if c.wantErr {
				asst.Error(err)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, got)
		})
	}
}
