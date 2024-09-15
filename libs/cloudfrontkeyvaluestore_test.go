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
