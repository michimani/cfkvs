package cli

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/michimani/cfkvs/internal/commands"
)

type CLI struct {
	commands.Globals

	Kvs  commands.KvsCmd  `cmd:"" help:"KeyValueStore operations."`
	Item commands.ItemCmd `cmd:"item" help:"Items in specific KeyValueStore."`
}

var (
	version       string
	revision      string
	versionString = fmt.Sprintf("%s-%s", version, revision)
)

func Run(ctx context.Context) error {
	var cli CLI
	var err error
	if cli, err = NewCLI(ctx); err != nil {
		return err
	}

	kctx := kong.Parse(&cli,
		kong.Name("cfkvs"),
		kong.Description("A simple cli tool to manage CloudFront Key Value Stores."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": versionString,
		})
	err = kctx.Run(&cli.Globals)
	kctx.FatalIfErrorf(err)

	return nil
}

func NewCLI(ctx context.Context) (CLI, error) {
	cli := CLI{
		Globals: commands.Globals{
			Version: commands.VersionFlag(versionString),
		},
	}

	return cli, nil
}
