package commands_test

import (
	"testing"

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
	}{}

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
	}{}

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
	}{}

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
	}{}

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
		wantError bool
	}{}

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
