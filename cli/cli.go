package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/michimani/cfkvs/internal/commands"
	"github.com/michimani/cfkvs/libs"
)

type CLI struct {
	commands.Globals

	KVS  commands.KVSCmd  `cmd:"" help:"KeyValueStore operations."`
	Item commands.ItemCmd `cmd:"item" help:"Items in specific KeyValueStore."`
}

var (
	version       string
	revision      string
	versionString = fmt.Sprintf("v%s-%s", version, revision)
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

	if err := setClient(ctx, kctx.Args, &cli.Globals); err != nil {
		return err
	}

	err = kctx.Run(&cli.Globals)
	kctx.FatalIfErrorf(err)

	return nil
}

var needClientCommands = []string{
	"item",
	"kvs",
}

func setClient(ctx context.Context, args []string, globals *commands.Globals) error {
	if len(args) == 0 {
		return nil
	}

	needClient := false
	for _, cmd := range needClientCommands {
		if args[0] == cmd {
			needClient = true
			break
		}
	}

	if !needClient {
		return nil
	}

	s3, err := libs.NewS3Client(ctx)
	if err != nil {
		return err
	}
	globals.S3Client = s3

	cf, err := libs.NewCloudFrontClient(ctx)
	if err != nil {
		return err
	}
	globals.CloudFrontClient = cf

	kvs, err := libs.NewCloudFrontKeyValueStoreClient(ctx)
	if err != nil {
		return err
	}
	globals.CloudFrontKeyValueStoreClient = kvs

	return nil
}

func NewCLI(ctx context.Context) (CLI, error) {
	cli := CLI{
		Globals: commands.Globals{
			Version:      commands.VersionFlag(versionString),
			OutputTarget: os.Stdout,
		},
	}

	return cli, nil
}
