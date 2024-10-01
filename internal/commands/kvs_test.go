package commands_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
	"github.com/michimani/cfkvs/internal/commands"
	"github.com/michimani/cfkvs/libs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_ListKvsSubCmd_Run(t *testing.T) {
	cases := []struct {
		name      string
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		wantError bool
	}{
		{
			name: "ok",
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).Return(
					&cf.ListKeyValueStoresOutput{
						KeyValueStoreList: &cfTypes.KeyValueStoreList{
							Items: []cfTypes.KeyValueStore{},
						},
					}, nil)
				return m
			},
			wantError: false,
		},
		{
			name: "error: ListKeyValueStores returns error",
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return m
			},
			wantError: true,
		},
		{
			name: "error: ListKeyValueStores invalid output",
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).Return(
					&cf.ListKeyValueStoresOutput{
						KeyValueStoreList: nil,
					}, nil)
				return m
			},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			ctrl := gomock.NewController(tt)
			cfcMock := c.cfcMock(ctrl)
			globals := &commands.Globals{
				CloudFrontClient: cfcMock,
			}

			err := (&commands.ListKvsSubCmd{}).Run(globals)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)
		})
	}
}

func Test_CreateSubCmd_Ruh(t *testing.T) {
	cases := []struct {
		name      string
		cmd       *commands.CreateSubCmd
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		wantError bool
	}{
		{
			name: "ok",
			cmd:  &commands.CreateSubCmd{},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.CreateKeyValueStoreOutput{
						KeyValueStore: &cfTypes.KeyValueStore{
							Id:      aws.String("id"),
							Name:    aws.String("name"),
							Comment: aws.String("comment"),
							Status:  aws.String("status"),
							ARN:     aws.String("arn"),
						},
					}, nil)
				return m
			},
			wantError: false,
		},
		{
			name: "ok: with source",
			cmd: &commands.CreateSubCmd{
				Bucket:    "bucket",
				ObjectKey: "object-key",
			},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.CreateKeyValueStoreOutput{
						KeyValueStore: &cfTypes.KeyValueStore{
							Id:      aws.String("id"),
							Name:    aws.String("name"),
							Comment: aws.String("comment"),
							Status:  aws.String("status"),
							ARN:     aws.String("arn"),
						},
					}, nil)
				return m
			},
			wantError: false,
		},
		{
			name: "error: object key is empty",
			cmd: &commands.CreateSubCmd{
				Bucket: "bucket",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			wantError: true,
		},
		{
			name: "error: CreateKeyValueStore returns error",
			cmd:  &commands.CreateSubCmd{},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return m
			},
			wantError: true,
		},
		{
			name: "error: CreateKeyValueStore invalid output",
			cmd:  &commands.CreateSubCmd{},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&cf.CreateKeyValueStoreOutput{
						KeyValueStore: nil,
					}, nil)
				return m
			},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			ctrl := gomock.NewController(tt)
			cfcMock := c.cfcMock(ctrl)
			globals := &commands.Globals{
				CloudFrontClient: cfcMock,
			}

			err := c.cmd.Run(globals)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)
		})
	}
}

func Test_InfoSubCmd_Run(t *testing.T) {
	cases := []struct {
		name      string
		cmd       *commands.InfoSubCmd
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		kvscMock  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient
		wantError bool
	}{
		{
			name: "ok",
			cmd:  &commands.InfoSubCmd{Name: "name"},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.DescribeKeyValueStoreOutput{
						KeyValueStore: &cfTypes.KeyValueStore{
							Id:      aws.String("id"),
							Name:    aws.String("name"),
							Comment: aws.String("comment"),
							Status:  aws.String("status"),
							ARN:     aws.String("arn"),
						},
					}, nil)
				return m
			},
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&kvs.DescribeKeyValueStoreOutput{
						KvsARN: aws.String("kvs-arn"),
					}, nil)
				return m
			},
			wantError: false,
		},
		{
			name: "error: cloudfront.DescribeKeyValueStore returns error",
			cmd:  &commands.InfoSubCmd{Name: "name"},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return m
			},
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: cloudfront.DescribeKeyValueStore invalid output",
			cmd:  &commands.InfoSubCmd{Name: "name"},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.DescribeKeyValueStoreOutput{
						KeyValueStore: nil,
					}, nil)
				return m
			},
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: cloudfrontkeyvaluestore.DescribeKeyValueStore returns error",
			cmd:  &commands.InfoSubCmd{Name: "name"},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.DescribeKeyValueStoreOutput{
						KeyValueStore: &cfTypes.KeyValueStore{
							Id:      aws.String("id"),
							Name:    aws.String("name"),
							Comment: aws.String("comment"),
							Status:  aws.String("status"),
							ARN:     aws.String("arn"),
						},
					}, nil)
				return m
			},
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return m
			},
			wantError: true,
		},
		{
			name: "error: cloudfrontkeyvaluestore.DescribeKeyValueStore invalid output",
			cmd:  &commands.InfoSubCmd{Name: "name"},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				m := libs.NewMockCloudFrontClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.DescribeKeyValueStoreOutput{
						KeyValueStore: &cfTypes.KeyValueStore{
							Id:      aws.String("id"),
							Name:    aws.String("name"),
							Comment: aws.String("comment"),
							Status:  aws.String("status"),
							ARN:     aws.String("arn"),
						},
					}, nil)
				return m
			},
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).Return(nil, nil)
				return m
			},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			ctrl := gomock.NewController(tt)
			cfcMock := c.cfcMock(ctrl)
			kvscMock := c.kvscMock(ctrl)
			globals := &commands.Globals{
				CloudFrontClient:              cfcMock,
				CloudFrontKeyValueStoreClient: kvscMock,
			}

			err := c.cmd.Run(globals)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)
		})
	}
}

func Test_SyncSubCmd_Run(t *testing.T) {
	cases := []struct {
		name      string
		cmd       *commands.SyncSubCmd
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		kvscMock  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient
		s3cMock   func(ctrl *gomock.Controller) *libs.MockS3Client
		wantError bool
	}{
		{
			name: "ok: from file",
			cmd: &commands.SyncSubCmd{
				Name: "kvs-name",
				File: "../../testdata/valid.json",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).
					Return(&kvs.ListKeysOutput{
						Items: []kvsTypes.ListKeysResponseListItem{},
					}, nil)
				return m
			},
			s3cMock: func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
		},
		{
			name: "ok: from S3 Object",
			cmd: &commands.SyncSubCmd{
				Name:      "kvs-name",
				Bucket:    "bucket",
				ObjectKey: "object-key",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).
					Return(&kvs.ListKeysOutput{
						Items: []kvsTypes.ListKeysResponseListItem{},
					}, nil)
				return m
			},
			s3cMock: noErrorMockS3Client,
		},
		{
			name: "ok: with yes",
			cmd: &commands.SyncSubCmd{
				Name: "kvs-name",
				File: "../../testdata/valid.json",
				Yes:  true,
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).
					Return(&kvs.ListKeysOutput{
						Items: []kvsTypes.ListKeysResponseListItem{},
					}, nil)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(&kvs.DescribeKeyValueStoreOutput{ETag: aws.String("etag")}, nil)
				m.EXPECT().UpdateKeys(gomock.Any(), gomock.Any()).Return(&kvs.UpdateKeysOutput{
					ItemCount:        aws.Int32(1),
					TotalSizeInBytes: aws.Int64(1024),
				}, nil)
				return m
			},
			s3cMock: func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
		},
		{
			name: "error: Name is empty",
			cmd: &commands.SyncSubCmd{
				Name: "",
				File: "../../testdata/valid.json",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			s3cMock:   func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
			wantError: true,
		},
		{
			name: "error: bucket is empty",
			cmd: &commands.SyncSubCmd{
				Name:      "kvs-name",
				ObjectKey: "object-key",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			s3cMock:   func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
			wantError: true,
		},
		{
			name: "error: object-key is empty",
			cmd: &commands.SyncSubCmd{
				Name:   "kvs-name",
				Bucket: "bucket",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			s3cMock:   func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
			wantError: true,
		},
		{
			name: "error: getKvsArn returns error",
			cmd: &commands.SyncSubCmd{
				Name: "kvs-name",
				File: "../../testdata/valid.json",
			},
			cfcMock:   errorMockCloudFrontClient,
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			s3cMock:   func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
			wantError: true,
		},
		{
			name: "error: libs.ListItems returns error",
			cmd: &commands.SyncSubCmd{
				Name: "kvs-name",
				File: "../../testdata/valid.json",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return m
			},
			s3cMock:   func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
			wantError: true,
		},
		{
			name: "error: libs.GetKeyValueStoreDataFromFile returns error",
			cmd: &commands.SyncSubCmd{
				Name: "kvs-name",
				File: "../../testdata/notfound.json",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).
					Return(&kvs.ListKeysOutput{
						Items: []kvsTypes.ListKeysResponseListItem{},
					}, nil)
				return m
			},
			s3cMock:   func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
			wantError: true,
		},
		{
			name: "error: libs.GetKeyValueStoreData returns error",
			cmd: &commands.SyncSubCmd{
				Name:      "kvs-name",
				Bucket:    "bucket",
				ObjectKey: "object-key",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).
					Return(&kvs.ListKeysOutput{
						Items: []kvsTypes.ListKeysResponseListItem{},
					}, nil)
				return m
			},
			s3cMock:   errorMockS3Client,
			wantError: true,
		},
		{
			name: "error: libs.SyncItems returns error",
			cmd: &commands.SyncSubCmd{
				Name: "kvs-name",
				File: "../../testdata/valid.json",
				Yes:  true,
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).
					Return(&kvs.ListKeysOutput{
						Items: []kvsTypes.ListKeysResponseListItem{},
					}, nil)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(&kvs.DescribeKeyValueStoreOutput{ETag: aws.String("etag")}, nil)
				m.EXPECT().UpdateKeys(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return m
			},
			s3cMock:   func(ctrl *gomock.Controller) *libs.MockS3Client { return nil },
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			ctrl := gomock.NewController(tt)
			cfcMock := c.cfcMock(ctrl)
			kvscMock := c.kvscMock(ctrl)
			s3cMock := c.s3cMock(ctrl)
			globals := &commands.Globals{
				CloudFrontClient:              cfcMock,
				CloudFrontKeyValueStoreClient: kvscMock,
				S3Client:                      s3cMock,
			}

			err := c.cmd.Run(globals)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)
		})
	}
}
