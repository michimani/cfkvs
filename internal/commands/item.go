package commands

import (
	"context"
	"errors"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/libs"
	"github.com/michimani/cfkvs/types"
)

type ItemCmd struct {
	List   ListItemsSubCmd `cmd:"" help:"List items in the key value store."`
	Get    GetSubCmd       `cmd:"" help:"Get an item in the key value store."`
	Put    PutSubCmd       `cmd:"" help:"Put an item in the key value store."`
	Delete DeleteSubCmd    `cmd:"" help:"Delete an item in the key value store."`
}

type ListItemsSubCmd struct {
	KvsName string `name:"kvs-name" help:"Name of the key value store." required:""`
}

type GetSubCmd struct {
	KvsName string `name:"kvs-name" help:"Name of the key value store." required:""`
	Key     string `name:"key" help:"Key of the item to get." required:""`
}

type PutSubCmd struct {
	KvsName string `name:"kvs-name" help:"Name of the key value store." required:""`
	Key     string `name:"key" help:"Key of the item to put." required:""`
	Value   string `name:"value" help:"Value of the item to put." required:""`
}

type DeleteSubCmd struct {
	KvsName string `name:"kvs-name" help:"Name of the key value store." required:""`
	Key     string `name:"key" help:"Key of the item to delete." required:""`
}

func getKvsArn(ctx context.Context, cfc libs.CloudFrontClient, kvsName string) (string, error) {
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
	kvsARN, err := getKvsArn(ctx, globals.CloudFrontClient, c.KvsName)
	if err != nil {
		return err
	}

	out, err := libs.ListItems(ctx, globals.CloudFrontKeyValueStoreClient, kvsARN)
	if err != nil {
		return err
	}

	itemList := types.ItemList{}
	if err := itemList.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&itemList, globals.Output, globals.OutputTarget); err != nil {
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
	kvsARN, err := getKvsArn(ctx, globals.CloudFrontClient, c.KvsName)
	if err != nil {
		return err
	}

	out, err := libs.GetItem(ctx, globals.CloudFrontKeyValueStoreClient, kvsARN, c.Key)
	if err != nil {
		return err
	}

	item := types.Item{}
	if err := item.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&item, globals.Output, globals.OutputTarget); err != nil {
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
	kvsARN, err := getKvsArn(ctx, globals.CloudFrontClient, c.KvsName)
	if err != nil {
		return err
	}

	out, err := libs.PutItem(ctx, globals.CloudFrontKeyValueStoreClient, kvsARN, c.Key, c.Value)
	if err != nil {
		return err
	}

	kvsSimple := types.KVSSimple{}
	if err := kvsSimple.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvsSimple, globals.Output, globals.OutputTarget); err != nil {
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
	kvsARN, err := getKvsArn(ctx, globals.CloudFrontClient, c.KvsName)
	if err != nil {
		return err
	}

	out, err := libs.DeleteItem(ctx, globals.CloudFrontKeyValueStoreClient, kvsARN, c.Key)
	if err != nil {
		return err
	}

	kvsSimple := types.KVSSimple{}
	if err := kvsSimple.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvsSimple, globals.Output, globals.OutputTarget); err != nil {
		return err
	}

	return nil
}
