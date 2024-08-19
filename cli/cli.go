package cli

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/michimani/cfkvs/internal/commands"
)

type Globals struct {
	Debug   bool        `short:"D" name:"debug" help:"Enable debug mode."`
	Version VersionFlag `name:"version" help:"Print version information and quit"`
}

type CLI struct {
	Globals

	Kvs  commands.KvsCmd  `cmd:"" help:"KeyValueStore operations."`
	Item commands.ItemCmd `cmd:"item" help:"Items in specific KeyValueStore."`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

func Run(ctx context.Context) error {
	var cli *CLI
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
			"version": "0.0.1",
		})
	err = kctx.Run(&cli.Globals)
	kctx.FatalIfErrorf(err)

	return nil
}

func NewCLI(ctx context.Context) (*CLI, error) {
	cli := CLI{
		Globals: Globals{
			Version: VersionFlag("0.1.1"),
		},
	}

	return &cli, nil
}
