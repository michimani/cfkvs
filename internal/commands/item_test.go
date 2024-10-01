package commands_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michimani/cfkvs/internal/commands"
	"github.com/michimani/cfkvs/libs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_ListItemsSubCmd_Run(t *testing.T) {
	cases := []struct {
		name      string
		cmd       *commands.ListItemsSubCmd
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		kvscMock  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient
		wantError bool
	}{
		{
			name: "ok",
			cmd: &commands.ListItemsSubCmd{
				KvsName: "kvs-name",
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
		},
		{
			name: "error: kvsName is empty",
			cmd: &commands.ListItemsSubCmd{
				KvsName: "",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: getKvsArn returns error",
			cmd: &commands.ListItemsSubCmd{
				KvsName: "kvs-name",
			},
			cfcMock:   errorMockCloudFrontClient,
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: libs.ListItems returns error",
			cmd: &commands.ListItemsSubCmd{
				KvsName: "kvs-name",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().ListKeys(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
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

func Test_GetSubCmd_Run(t *testing.T) {
	cases := []struct {
		name      string
		cmd       *commands.GetSubCmd
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		kvscMock  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient
		wantError bool
	}{
		{
			name: "ok",
			cmd: &commands.GetSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().GetKey(gomock.Any(), gomock.Any()).
					Return(&kvs.GetKeyOutput{
						Key:   aws.String("key"),
						Value: aws.String("value"),
					}, nil)
				return m
			},
		},
		{
			name: "error: kvsName is empty",
			cmd: &commands.GetSubCmd{
				KvsName: "",
				Key:     "key",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: key is empty",
			cmd: &commands.GetSubCmd{
				KvsName: "kvs-name",
				Key:     "",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: getKvsArn returns error",
			cmd: &commands.GetSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
			},
			cfcMock:   errorMockCloudFrontClient,
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: libs.GetItem returns error",
			cmd: &commands.GetSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().GetKey(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
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

func Test_PutSubCmd_Run(t *testing.T) {
	cases := []struct {
		name      string
		cmd       *commands.PutSubCmd
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		kvscMock  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient
		wantError bool
	}{
		{
			name: "ok",
			cmd: &commands.PutSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
				Value:   "value",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(&kvs.DescribeKeyValueStoreOutput{
						ETag: aws.String("etag"),
					}, nil)
				m.EXPECT().PutKey(gomock.Any(), gomock.Any()).
					Return(&kvs.PutKeyOutput{
						ItemCount:        aws.Int32(1),
						TotalSizeInBytes: aws.Int64(1024),
					}, nil)
				return m
			},
		},
		{
			name: "error: kvsName is empty",
			cmd: &commands.PutSubCmd{
				KvsName: "",
				Key:     "key",
				Value:   "value",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: key is empty",
			cmd: &commands.PutSubCmd{
				KvsName: "kvs-name",
				Key:     "",
				Value:   "value",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: value is empty",
			cmd: &commands.PutSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
				Value:   "",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: getKvsArn returns error",
			cmd: &commands.PutSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
				Value:   "value",
			},
			cfcMock:   errorMockCloudFrontClient,
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: libs.PutItem returns error (failed to get ETag)",
			cmd: &commands.PutSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
				Value:   "value",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
				return m
			},
			wantError: true,
		},
		{
			name: "error: libs.PutItem returns error (failed to put item)",
			cmd: &commands.PutSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
				Value:   "value",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(&kvs.DescribeKeyValueStoreOutput{
						ETag: aws.String("etag"),
					}, nil)
				m.EXPECT().PutKey(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
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

func Test_DeleteSubCmd_Run(t *testing.T) {
	cases := []struct {
		name      string
		cmd       *commands.DeleteSubCmd
		cfcMock   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient
		kvscMock  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient
		wantError bool
	}{
		{
			name: "ok",
			cmd: &commands.DeleteSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(&kvs.DescribeKeyValueStoreOutput{
						ETag: aws.String("etag"),
					}, nil)
				m.EXPECT().DeleteKey(gomock.Any(), gomock.Any()).
					Return(&kvs.DeleteKeyOutput{
						ItemCount:        aws.Int32(1),
						TotalSizeInBytes: aws.Int64(1024),
					}, nil)
				return m
			},
		},
		{
			name: "error: kvsName is empty",
			cmd: &commands.DeleteSubCmd{
				KvsName: "",
				Key:     "key",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: key is empty",
			cmd: &commands.DeleteSubCmd{
				KvsName: "kvs-name",
				Key:     "",
			},
			cfcMock:   func(ctrl *gomock.Controller) *libs.MockCloudFrontClient { return nil },
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: getKvsArn returns error",
			cmd: &commands.DeleteSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
			},
			cfcMock:   errorMockCloudFrontClient,
			kvscMock:  func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient { return nil },
			wantError: true,
		},
		{
			name: "error: libs.DeleteItem returns error (failed to get ETag)",
			cmd: &commands.DeleteSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
				return m
			},
			wantError: true,
		},
		{
			name: "error: libs.DeleteItem returns error (failed to delete item)",
			cmd: &commands.DeleteSubCmd{
				KvsName: "kvs-name",
				Key:     "key",
			},
			cfcMock: noErrorMockCloudFrontClient,
			kvscMock: func(ctrl *gomock.Controller) *libs.MockCloudFrontKeyValueStoreClient {
				m := libs.NewMockCloudFrontKeyValueStoreClient(ctrl)
				m.EXPECT().DescribeKeyValueStore(gomock.Any(), gomock.Any()).
					Return(&kvs.DescribeKeyValueStoreOutput{
						ETag: aws.String("etag"),
					}, nil)
				m.EXPECT().DeleteKey(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
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

func errorMockCloudFrontClient(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
	m := libs.NewMockCloudFrontClient(ctrl)
	m.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).
		Return(&cf.ListKeyValueStoresOutput{
			KeyValueStoreList: &cfTypes.KeyValueStoreList{
				Items: []cfTypes.KeyValueStore{},
			},
		}, nil)
	return m
}

func noErrorMockCloudFrontClient(ctrl *gomock.Controller) *libs.MockCloudFrontClient {
	m := libs.NewMockCloudFrontClient(ctrl)
	m.EXPECT().ListKeyValueStores(gomock.Any(), gomock.Any()).
		Return(&cf.ListKeyValueStoresOutput{
			KeyValueStoreList: &cfTypes.KeyValueStoreList{
				Items: []cfTypes.KeyValueStore{
					{Name: aws.String("kvs-name"), ARN: aws.String("kvs-arn")},
				},
			},
		}, nil)
	return m
}

func errorMockS3Client(ctrl *gomock.Controller) *libs.MockS3Client {
	m := libs.NewMockS3Client(ctrl)
	m.EXPECT().GetObject(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("error"))
	return m
}

func noErrorMockS3Client(ctrl *gomock.Controller) *libs.MockS3Client {
	m := libs.NewMockS3Client(ctrl)
	m.EXPECT().GetObject(gomock.Any(), gomock.Any()).
		Return(&s3.GetObjectOutput{
			Body: io.NopCloser(strings.NewReader(`{"data": [{"key":"k", "value":"v"}]}`)),
		}, nil)
	return m
}
