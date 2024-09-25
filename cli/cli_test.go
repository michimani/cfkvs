package cli_test

import (
	"context"
	"testing"

	"github.com/michimani/cfkvs/cli"
	"github.com/michimani/cfkvs/internal/commands"
	"github.com/stretchr/testify/assert"
)

func Test_NewCLI(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		version       string
		revision      string
		expectVersion commands.VersionFlag
		wantError     bool
	}{
		{
			name:          "ok",
			version:       "1.0.0",
			revision:      "abcdefg",
			expectVersion: commands.VersionFlag("v1.0.0-abcdefg"),
			ctx:           context.TODO(),
		},
		{
			name:          "ok: empty version",
			version:       "",
			revision:      "",
			expectVersion: commands.VersionFlag("v-"),
			ctx:           context.TODO(),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			cli.SetVersion(c.version)
			cli.SetRevision(c.revision)
			cli.SetVersionString()

			ac, err := cli.NewCLI(c.ctx)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expectVersion, ac.Version)
		})
	}
}

func Test_setClient(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		globals   *commands.Globals
		envs      map[string]string
		wantSet   bool
		wantError bool
	}{
		{
			name:    "ok: want set client for item command",
			args:    []string{"item"},
			globals: &commands.Globals{},
			envs: map[string]string{
				"AWS_ACCESS_KEY_ID":     "dummy_key_id",
				"AWS_SECRET_ACCESS_KEY": "dummy_secret_key",
				"AWS_SESSION_TOKEN":     "dummy_session_token",
				"AWS_REGION":            "ap-northeast-1",
			},
			wantSet: true,
		},
		{
			name:    "ok: want set client for kvs command",
			args:    []string{"kvs"},
			globals: &commands.Globals{},
			envs: map[string]string{
				"AWS_ACCESS_KEY_ID":     "dummy_key_id",
				"AWS_SECRET_ACCESS_KEY": "dummy_secret_key",
				"AWS_SESSION_TOKEN":     "dummy_session_token",
				"AWS_REGION":            "ap-northeast-1",
			},
			wantSet: true,
		},
		{
			name:    "ok: not want set client for other command",
			args:    []string{"other"},
			globals: &commands.Globals{},
			wantSet: false,
		},
		{
			name:    "error: failed to new client",
			args:    []string{"item"},
			globals: &commands.Globals{},
			envs: map[string]string{
				"AWS_PROFILE": "invalid",
				"AWS_REGION":  "ap-northeast-1",
			},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			for k, v := range c.envs {
				tt.Setenv(k, v)
			}

			err := cli.Exported_setClient(context.TODO(), c.args, c.globals)
			if c.wantError {
				asst.Error(err)
				return
			}

			asst.NoError(err)
			if c.wantSet {
				asst.NotNil(c.globals.S3Client)
				asst.NotNil(c.globals.CloudFrontClient)
				asst.NotNil(c.globals.CloudFrontKeyValueStoreClient)
			} else {
				asst.Nil(c.globals.S3Client)
				asst.Nil(c.globals.CloudFrontClient)
				asst.Nil(c.globals.CloudFrontKeyValueStoreClient)
			}
		})
	}
}
