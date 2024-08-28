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

func getKvsArn(ctx context.Context, kvsName string) (string, error) {
	cfc, err := libs.NewCloudFrontClient(ctx)
	if err != nil {
		return "", err
	}

	kvsARN, err := libs.GetKeyValueStoreArn(ctx, cfc, kvsName)
	if err != nil {
		return "", err
	}

	return kvsARN, nil
}

func (c *ListItemsSubCmd) Run(globals *Globals) error {
	if c.KvsName == "" {
		return errors.New("kvs-name is required")
	}

	ctx := context.TODO()
	kvsc, err := libs.NewCloudFrontKeyValueStoreClient(ctx)
	if err != nil {
		return err
	}

	kvsARN, err := getKvsArn(ctx, c.KvsName)
	if err != nil {
		return err
	}

	itemList, err := libs.ListItems(ctx, kvsc, kvsARN)
	if err != nil {
		return err
	}

	if err := output.Render(itemList, globals.Output); err != nil {
		return err
	}

	return nil
}

func (c *GetSubCmd) Run(globals *Globals) error {
	if c.KvsName == "" {
		return errors.New("kvs-name is required")
	}
	if c.Key == "" {
		return errors.New("key is required")
	}

	ctx := context.TODO()
	kvsc, err := libs.NewCloudFrontKeyValueStoreClient(ctx)
	if err != nil {
		return err
	}

	kvsARN, err := getKvsArn(ctx, c.KvsName)
	if err != nil {
		return err
	}

	item, err := libs.GetItem(ctx, kvsc, kvsARN, c.Key)
	if err != nil {
		return err
	}

	if err := output.Render(item, globals.Output); err != nil {
		return err
	}

	return nil
}

func (c *PutSubCmd) Run(globals *Globals) error {
	if c.KvsName == "" {
		return errors.New("kvs-name is required")
	}
	if c.Key == "" {
		return errors.New("key is required")
	}
	if c.Value == "" {
		return errors.New("value is required")
	}

	ctx := context.TODO()
	kvsc, err := libs.NewCloudFrontKeyValueStoreClient(ctx)
	if err != nil {
		return err
	}

	kvsARN, err := getKvsArn(ctx, c.KvsName)
	if err != nil {
		return err
	}

	res, err := libs.PutItem(ctx, kvsc, kvsARN, c.Key, c.Value)
	if err != nil {
		return err
	}

	if err := output.Render(res, globals.Output); err != nil {
		return err
	}

	return nil
}

func (c *DeleteSubCmd) Run(globals *Globals) error {
	if c.KvsName == "" {
		return errors.New("kvs-name is required")
	}
	if c.Key == "" {
		return errors.New("key is required")
	}

	ctx := context.TODO()
	kvsc, err := libs.NewCloudFrontKeyValueStoreClient(ctx)
	if err != nil {
		return err
	}

	kvsARN, err := getKvsArn(ctx, c.KvsName)
	if err != nil {
		return err
	}

	res, err := libs.DeleteItem(ctx, kvsc, kvsARN, c.Key)
	if err != nil {
		return err
	}

	if err := output.Render(res, globals.Output); err != nil {
		return err
	}

	return nil
}
