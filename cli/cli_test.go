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
