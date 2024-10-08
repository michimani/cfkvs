package commands

import (
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/libs"
)

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type Globals struct {
	Debug   bool              `short:"D" name:"debug" help:"Enable debug mode."`
	Output  output.OutputType `name:"output" help:"Output format. One of: json, table." default:"table"`
	Version VersionFlag       `name:"version" help:"Print version information and quit"`

	S3Client                      libs.S3Client                      `kong:"-"`
	CloudFrontClient              libs.CloudFrontClient              `kong:"-"`
	CloudFrontKeyValueStoreClient libs.CloudFrontKeyValueStoreClient `kong:"-"`
	OutputTarget                  io.Writer                          `kong:"-"`
}
