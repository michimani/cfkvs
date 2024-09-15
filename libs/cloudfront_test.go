package libs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/michimani/cfkvs/libs"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_NewCloudFrontClient(t *testing.T) {
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

			sc, err := libs.NewCloudFrontClient(c.ctx)
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

func Test_GetKeyValueStoreArn(t *testing.T) {
	cases := []struct {
		name      string
		clientOut struct {
			ListKeyValueStoresOutput *cloudfront.ListKeyValueStoresOutput
			Error                    error
		}
		kvsName string
		wantArn string
		wantErr bool
	}{
		{
			name: "success",
			clientOut: struct {
				ListKeyValueStoresOutput *cloudfront.ListKeyValueStoresOutput
				Error                    error
			}{
				ListKeyValueStoresOutput: &cloudfront.ListKeyValueStoresOutput{
					KeyValueStoreList: &cfTypes.KeyValueStoreList{
						Items: []cfTypes.KeyValueStore{
							{
								Name: aws.String("kvs_name"),
								ARN:  aws.String("arn"),
							},
						},
					},
				},
				Error: nil,
			},
			kvsName: "kvs_name",
			wantArn: "arn",
			wantErr: false,
		},
		{
			name: "failed to list key value stores",
			clientOut: struct {
				ListKeyValueStoresOutput *cloudfront.ListKeyValueStoresOutput
				Error                    error
			}{
				ListKeyValueStoresOutput: nil,
				Error:                    errors.New("failed to list key value stores"),
			},
			kvsName: "kvs_name",
			wantArn: "",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			ctx := context.Background()

			ctrl := gomock.NewController(tt)
			m := libs.NewMockCloudFrontClient(ctrl)
			m.EXPECT().
				ListKeyValueStores(ctx, &cloudfront.ListKeyValueStoresInput{}).
				Return(c.clientOut.ListKeyValueStoresOutput, c.clientOut.Error)

			arn, err := libs.GetKeyValueStoreArn(ctx, m, c.kvsName)
			if c.wantErr {
				asst.Error(err)
				asst.Empty(arn)
				return
			}

			asst.NoError(err)
			asst.Equal(c.wantArn, arn)
		})
	}
}

func Test_ListKvs(t *testing.T) {
	cases := []struct {
		name      string
		clientOut struct {
			ListKeyValueStoresOutput *cloudfront.ListKeyValueStoresOutput
			Error                    error
		}
		wantErr bool
	}{
		{
			name: "success",
			clientOut: struct {
				ListKeyValueStoresOutput *cloudfront.ListKeyValueStoresOutput
				Error                    error
			}{
				ListKeyValueStoresOutput: &cloudfront.ListKeyValueStoresOutput{},
				Error:                    nil,
			},
			wantErr: false,
		},
		{
			name: "failed to list key value stores",
			clientOut: struct {
				ListKeyValueStoresOutput *cloudfront.ListKeyValueStoresOutput
				Error                    error
			}{
				ListKeyValueStoresOutput: nil,
				Error:                    errors.New("failed to list key value stores"),
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			ctx := context.Background()

			ctrl := gomock.NewController(tt)
			m := libs.NewMockCloudFrontClient(ctrl)
			m.EXPECT().
				ListKeyValueStores(ctx, &cloudfront.ListKeyValueStoresInput{}).
				Return(c.clientOut.ListKeyValueStoresOutput, c.clientOut.Error)

			out, err := libs.ListKvs(ctx, m)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(out)
				return
			}

			asst.NoError(err)
			asst.NotNil(out)
		})
	}
}

func Test_KVSImportSourceS3_ARN(t *testing.T) {
	cases := []struct {
		name string
		s3   libs.KVSImportSourceS3
		want string
	}{
		{
			name: "success",
			s3:   libs.KVSImportSourceS3{Bucket: "bucket", Key: "key"},
			want: "arn:aws:s3:::bucket/key",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			asst.Equal(c.want, c.s3.ARN())
		})
	}
}

func Test_KVSImportSourceS3_Type(t *testing.T) {
	cases := []struct {
		name string
		s3   libs.KVSImportSourceS3
		want cfTypes.ImportSourceType
	}{
		{
			name: "success",
			s3:   libs.KVSImportSourceS3{Bucket: "bucket", Key: "key"},
			want: cfTypes.ImportSourceTypeS3,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			asst.Equal(c.want, c.s3.Type())
		})
	}
}

func Test_CreateKvs(t *testing.T) {
	cases := []struct {
		name      string
		clientOut struct {
			CreateKeyValueStoreOutput *cloudfront.CreateKeyValueStoreOutput
			Error                     error
		}
		kvsName string
		comment string
		source  libs.KVSImportSource
		wantErr bool
	}{
		{
			name: "success",
			clientOut: struct {
				CreateKeyValueStoreOutput *cloudfront.CreateKeyValueStoreOutput
				Error                     error
			}{
				CreateKeyValueStoreOutput: &cloudfront.CreateKeyValueStoreOutput{
					KeyValueStore: &cfTypes.KeyValueStore{
						Id:      aws.String("id"),
						Name:    aws.String("kvs_name"),
						Comment: aws.String("kvs_comment"),
						Status:  aws.String("status"),
						ARN:     aws.String("arn"),
					},
				},
				Error: nil,
			},
			kvsName: "kvs_name",
			comment: "kvs_comment",
			wantErr: false,
		},
		{
			name: "success with import source",
			clientOut: struct {
				CreateKeyValueStoreOutput *cloudfront.CreateKeyValueStoreOutput
				Error                     error
			}{
				CreateKeyValueStoreOutput: &cloudfront.CreateKeyValueStoreOutput{
					KeyValueStore: &cfTypes.KeyValueStore{
						Id:      aws.String("id"),
						Name:    aws.String("kvs_name"),
						Comment: aws.String("kvs_comment"),
						Status:  aws.String("status"),
						ARN:     aws.String("arn"),
					},
				},
				Error: nil,
			},
			kvsName: "kvs_name",
			comment: "kvs_comment",
			source:  libs.KVSImportSourceS3{Bucket: "bucket", Key: "key"},
			wantErr: false,
		},
		{
			name: "failed to create key value store",
			clientOut: struct {
				CreateKeyValueStoreOutput *cloudfront.CreateKeyValueStoreOutput
				Error                     error
			}{
				CreateKeyValueStoreOutput: nil,
				Error:                     errors.New("failed to create key value store"),
			},
			kvsName: "kvs_name",
			comment: "kvs_comment",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			ctx := context.Background()

			ctrl := gomock.NewController(tt)
			m := libs.NewMockCloudFrontClient(ctrl)

			if c.source != nil {
				is := new(cfTypes.ImportSource)
				is.SourceARN = aws.String(c.source.ARN())
				is.SourceType = c.source.Type()

				m.EXPECT().
					CreateKeyValueStore(ctx, &cloudfront.CreateKeyValueStoreInput{
						Name:         &c.kvsName,
						Comment:      &c.comment,
						ImportSource: is,
					}).
					Return(c.clientOut.CreateKeyValueStoreOutput, c.clientOut.Error)
			} else {
				m.EXPECT().
					CreateKeyValueStore(ctx, &cloudfront.CreateKeyValueStoreInput{
						Name:    &c.kvsName,
						Comment: &c.comment,
					}).
					Return(c.clientOut.CreateKeyValueStoreOutput, c.clientOut.Error)
			}

			out, err := libs.CreateKvs(ctx, m, c.kvsName, c.comment, c.source)
			if c.wantErr {
				asst.Error(err)
				asst.Nil(out)
				return
			}

			asst.NoError(err)
			asst.NotNil(out)
		})
	}
}
