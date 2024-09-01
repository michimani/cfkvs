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
	Sync   SyncSubCmd      `cmd:"" help:"Sync items in the key value store with S3 object."`
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

type SyncSubCmd struct {
	KvsName   string `name:"kvs-name" help:"Name of the key value store." required:""`
	Bucket    string `name:"bucket" help:"S3 bucket name to sync key value store" required:""`
	ObjectKey string `name:"object-key" help:"S3 object key to sync key value store" required:""`
	Delete    bool   `name:"delete" help:"Delete items that are not in the S3 object."`
	Yes       bool   `name:"yes" short:"y" help:"Execute sync. If not specified, only show the items to be synced."`
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

	out, err := libs.ListItems(ctx, kvsc, kvsARN)
	if err != nil {
		return err
	}

	itemList := types.ItemList{}
	if err := itemList.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&itemList, globals.Output); err != nil {
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

	out, err := libs.GetItem(ctx, kvsc, kvsARN, c.Key)
	if err != nil {
		return err
	}

	item := types.Item{}
	if err := item.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&item, globals.Output); err != nil {
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

	out, err := libs.PutItem(ctx, kvsc, kvsARN, c.Key, c.Value)
	if err != nil {
		return err
	}

	kvsSimple := types.KVSSimple{}
	if err := kvsSimple.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvsSimple, globals.Output); err != nil {
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

	out, err := libs.DeleteItem(ctx, kvsc, kvsARN, c.Key)
	if err != nil {
		return err
	}

	kvsSimple := types.KVSSimple{}
	if err := kvsSimple.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvsSimple, globals.Output); err != nil {
		return err
	}

	return nil
}

func (c *SyncSubCmd) Run(globals *Globals) error {
	if c.KvsName == "" {
		return errors.New("kvs-name is required")
	}
	if c.Bucket == "" {
		return errors.New("bucket is required")
	}
	if c.ObjectKey == "" {
		return errors.New("object-key is required")
	}

	ctx := context.TODO()
	kvsc, err := libs.NewCloudFrontKeyValueStoreClient(ctx)
	if err != nil {
		return err
	}

	s3c, err := libs.NewS3Client(ctx)
	if err != nil {
		return err
	}

	kvsARN, err := getKvsArn(ctx, c.KvsName)
	if err != nil {
		return err
	}

	// get before items
	beforeOut, err := libs.ListItems(ctx, kvsc, kvsARN)
	if err != nil {
		return err
	}
	before := types.ItemList{}
	if err := before.Parse(beforeOut); err != nil {
		return err
	}

	// get after items from S3
	afterData, err := libs.GetKeyValueStoreData(ctx, s3c, c.Bucket, c.ObjectKey)
	if err != nil {
		return err
	}

	// show diff
	diff := before.Diff(afterData.ToItemList(), c.Delete)
	if err := output.Render(diff, output.OutputTypeTable); err != nil {
		return err
	}

	if !c.Yes {
		// TODO: sync
		return nil
	}

	return nil
}
