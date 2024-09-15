package libs_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
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

func Test_NewNewCloudFrontKeyValueStoreClient(t *testing.T) {
	cases := []struct {
		name    string
		ctx     context.Context
		envs    map[string]string
		wantErr bool
	}{
		{
			name: "ok",
			ctx:  context.Background(),
			envs: map[string]string{
				"AWS_ACCESS_KEY_ID":     "dummy_key_id",
				"AWS_SECRET_ACCESS_KEY": "dummy_secret_key",
				"AWS_SESSION_TOKEN":     "dummy_session_token",
				"AWS_REGION":            "ap-northeast-1",
			},
			wantErr: false,
		},
		{
			name: "failed to load default config",
			ctx:  context.Background(),
			envs: map[string]string{
				"AWS_PROFILE": "invalid",
				"AWS_REGION":  "ap-northeast-1",
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			for k, v := range c.envs {
				tt.Setenv(k, v)
			}

			asst := assert.New(tt)

			sc, err := libs.NewCloudFrontKeyValueStoreClient(c.ctx)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(sc)
				return
			}

			asst.NoError(err)
			asst.NotNil(sc)
		})
	}
}

func Test_ListItems(t *testing.T) {
	cases := []struct {
		name    string
		kvscOut struct {
			Out   *kvs.ListKeysOutput
			Error error
		}
		kvsARN  string
		expect  *kvs.ListKeysOutput
		wantErr bool
	}{
		{
			name: "ok",
			kvscOut: struct {
				Out   *kvs.ListKeysOutput
				Error error
			}{
				Out: &kvs.ListKeysOutput{
					Items: []kvsTypes.ListKeysResponseListItem{
						{Key: aws.String("key1"), Value: aws.String("value1")},
					},
				},
				Error: nil,
			},
			kvsARN: "dummy_arn",
			expect: &kvs.ListKeysOutput{
				Items: []kvsTypes.ListKeysResponseListItem{
					{Key: aws.String("key1"), Value: aws.String("value1")},
				},
			},
			wantErr: false,
		},
		{
			name: "failed to list keys",
			kvscOut: struct {
				Out   *kvs.ListKeysOutput
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
				ListKeys(gomock.Any(), &kvs.ListKeysInput{
					KvsARN: aws.String(c.kvsARN),
				}).
				Return(c.kvscOut.Out, c.kvscOut.Error)

			got, err := libs.ListItems(context.Background(), m, c.kvsARN)
			if c.wantErr {
				asst.Error(err)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, got)
		})
	}
}

func Test_GetItem(t *testing.T) {
	cases := []struct {
		name    string
		kvscOut struct {
			Out   *kvs.GetKeyOutput
			Error error
		}
		kvsARN  string
		key     string
		expect  *kvs.GetKeyOutput
		wantErr bool
	}{
		{
			name: "ok",
			kvscOut: struct {
				Out   *kvs.GetKeyOutput
				Error error
			}{
				Out: &kvs.GetKeyOutput{
					Key:   aws.String("key1"),
					Value: aws.String("value1"),
				},
				Error: nil,
			},
			kvsARN: "dummy_arn",
			key:    "key1",
			expect: &kvs.GetKeyOutput{
				Key:   aws.String("key1"),
				Value: aws.String("value1"),
			},
			wantErr: false,
		},
		{
			name: "failed to get key",
			kvscOut: struct {
				Out   *kvs.GetKeyOutput
				Error error
			}{
				Out:   nil,
				Error: assert.AnError,
			},
			kvsARN:  "dummy_arn",
			key:     "key1",
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
				GetKey(gomock.Any(), &kvs.GetKeyInput{
					KvsARN: aws.String(c.kvsARN),
					Key:    aws.String(c.key),
				}).
				Return(c.kvscOut.Out, c.kvscOut.Error)

			got, err := libs.GetItem(context.Background(), m, c.kvsARN, c.key)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(got)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, got)
		})
	}
}
