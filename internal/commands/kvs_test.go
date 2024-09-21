package commands_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
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
				cfcMock := libs.NewMockCloudFrontClient(ctrl)
				cfcMock.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).Return(
					&cf.ListKeyValueStoresOutput{
						KeyValueStoreList: &cfTypes.KeyValueStoreList{
							Items: []cfTypes.KeyValueStore{},
						},
					}, nil)
				return cfcMock
			},
			wantError: false,
		},
		{
			name: "error: ListKeyValueStores returns error",
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				cfcMock := libs.NewMockCloudFrontClient(ctrl)
				cfcMock.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return cfcMock
			},
			wantError: true,
		},
		{
			name: "error: ListKeyValueStores invalid output",
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				cfcMock := libs.NewMockCloudFrontClient(ctrl)
				cfcMock.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).Return(
					&cf.ListKeyValueStoresOutput{
						KeyValueStoreList: nil,
					}, nil)
				return cfcMock
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
				cfcMock := libs.NewMockCloudFrontClient(ctrl)
				cfcMock.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.CreateKeyValueStoreOutput{
						KeyValueStore: &cfTypes.KeyValueStore{
							Id:      aws.String("id"),
							Name:    aws.String("name"),
							Comment: aws.String("comment"),
							Status:  aws.String("status"),
							ARN:     aws.String("arn"),
						},
					}, nil)
				return cfcMock
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
				cfcMock := libs.NewMockCloudFrontClient(ctrl)
				cfcMock.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any()).Return(
					&cf.CreateKeyValueStoreOutput{
						KeyValueStore: &cfTypes.KeyValueStore{
							Id:      aws.String("id"),
							Name:    aws.String("name"),
							Comment: aws.String("comment"),
							Status:  aws.String("status"),
							ARN:     aws.String("arn"),
						},
					}, nil)
				return cfcMock
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
				cfcMock := libs.NewMockCloudFrontClient(ctrl)
				cfcMock.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return cfcMock
			},
			wantError: true,
		},
		{
			name: "error: CreateKeyValueStore invalid output",
			cmd:  &commands.CreateSubCmd{},
			cfcMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
				cfcMock := libs.NewMockCloudFrontClient(ctrl)
				cfcMock.EXPECT().CreateKeyValueStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&cf.CreateKeyValueStoreOutput{
						KeyValueStore: nil,
					}, nil)
				return cfcMock
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
