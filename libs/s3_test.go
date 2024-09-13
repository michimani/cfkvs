package libs_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michimani/cfkvs/libs"
	"github.com/michimani/cfkvs/types"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_NewS3Client(t *testing.T) {
	cases := []struct {
		name    string
		ctx     context.Context
		envs    map[string]string
		wantErr bool
	}{
		{
			name: "success",
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

			sc, err := libs.NewS3Client(c.ctx)
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

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}

func Test_GetKeyValueStoreData(t *testing.T) {
	cases := []struct {
		name      string
		ctx       context.Context
		clientOut struct {
			GetObjectOutput *s3.GetObjectOutput
			Error           error
		}
		bucket  string
		key     string
		expect  *types.KeyValueStoreData
		wantErr bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			clientOut: struct {
				GetObjectOutput *s3.GetObjectOutput
				Error           error
			}{
				GetObjectOutput: &s3.GetObjectOutput{
					Body: io.NopCloser(strings.NewReader(`{"data": [{"key":"k", "value":"v"}]}`)),
				},
				Error: nil,
			},
			bucket: "test-bucket",
			key:    "test-key",
			expect: &types.KeyValueStoreData{
				Data: &[]types.Item{
					{
						Key:   "k",
						Value: "v",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failed to get object",
			ctx:  context.Background(),
			clientOut: struct {
				GetObjectOutput *s3.GetObjectOutput
				Error           error
			}{
				GetObjectOutput: nil,
				Error:           errors.New("failed to get object"),
			},
			bucket:  "test-bucket",
			key:     "test-key",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "failed to read all",
			ctx:  context.Background(),
			clientOut: struct {
				GetObjectOutput *s3.GetObjectOutput
				Error           error
			}{
				GetObjectOutput: &s3.GetObjectOutput{
					Body: &errorReader{},
				},
				Error: nil,
			},
			bucket:  "test-bucket",
			key:     "test-key",
			expect:  nil,
			wantErr: true,
		},
		{
			name: "failed to unmarshal",
			ctx:  context.Background(),
			clientOut: struct {
				GetObjectOutput *s3.GetObjectOutput
				Error           error
			}{
				GetObjectOutput: &s3.GetObjectOutput{
					Body: io.NopCloser(strings.NewReader(`invalid`)),
				},
				Error: nil,
			},
			bucket:  "test-bucket",
			key:     "test-key",
			expect:  nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			ctrl := gomock.NewController(tt)
			m := libs.NewMockS3Client(ctrl)
			m.EXPECT().GetObject(c.ctx, &s3.GetObjectInput{
				Bucket: &c.bucket,
				Key:    &c.key,
			}).Return(c.clientOut.GetObjectOutput, c.clientOut.Error)

			kvsData, err := libs.GetKeyValueStoreData(c.ctx, m, c.bucket, c.key)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(kvsData)
				return
			}

			asst.NoError(err)
			asst.NotNil(kvsData)
		})
	}
}
