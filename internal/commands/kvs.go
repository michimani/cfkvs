package commands

import (
	"context"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/libs"
)

type KvsCmd struct {
	List   ListKvsSubCmd `cmd:"" help:"List key value stores in your account."`
	Create CreateSubCmd  `cmd:"" help:"Create a key value store."`
}

type ListKvsSubCmd struct{}

type CreateSubCmd struct {
	Name    string `arg:"" name:"name" help:"Name of the key value store." required:""`
	Comment string `arg:"" name:"comment" help:"Comment of the key value store."`
}

func (c *ListKvsSubCmd) Run(globals *Globals) error {
	cfc, err := libs.NewCloudFrontClient(context.TODO())
	if err != nil {
		return err
	}

	kvsList, err := libs.ListKvs(context.TODO(), cfc)
	if err != nil {
		return err
	}

	if err := output.RenderAsTable(kvsList); err != nil {
		return err
	}

	return nil
}
