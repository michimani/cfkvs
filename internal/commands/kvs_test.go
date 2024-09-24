package commands_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
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
