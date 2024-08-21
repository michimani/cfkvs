package commands

import (
	"context"
	"errors"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/libs"
)

type ItemCmd struct {
	List   ListItemsSubCmd `cmd:"" help:"List items in the key value store."`
	Get    GetSubCmd       `cmd:"" help:"Get an item in the key value store."`
	Put    PutSubCmd       `cmd:"" help:"Put an item in the key value store."`
	Delete DeleteSubCmd    `cmd:"" help:"Delete an item in the key value store."`
}

type ListItemsSubCmd struct {
	KvsName string `arg:"" name:"kvs-name" help:"Name of the key value store." required:""`
}

type GetSubCmd struct {
	KvsName string `arg:"" name:"kvs-name" help:"Name of the key value store." required:""`
	Key     string `arg:"" name:"key" help:"Key of the item to get." required:""`
}

type PutSubCmd struct {
	KvsName string `arg:"" name:"kvs-name" help:"Name of the key value store." required:""`
	Key     string `arg:"" name:"key" help:"Key of the item to put." required:""`
	Value   string `arg:"" name:"value" help:"Value of the item to put." required:""`
}

type DeleteSubCmd struct {
	KvsName string `arg:"" name:"kvs-name" help:"Name of the key value store." required:""`
	Key     string `arg:"" name:"key" help:"Key of the item to delete." required:""`
}

func (c *ListItemsSubCmd) Run(globals *Globals) error {
	if c.KvsName == "" {
		return errors.New("kvs-name is required")
	}

	ctx := context.TODO()
	cfc, err := libs.NewCloudFrontClient(ctx)
	if err != nil {
		return err
	}

	kvsc, err := libs.NewCloudFrontKeyValueStoreClient(ctx)
	if err != nil {
		return err
	}

	kvsARN, err := libs.GetKeyValueStoreArn(ctx, cfc, c.KvsName)
	if err != nil {
		return err
	}

	itemList, err := libs.ListItems(ctx, kvsc, kvsARN)
	if err != nil {
		return err
	}

	if err := output.RenderAsTable(itemList); err != nil {
		return err
	}

	return nil
}
