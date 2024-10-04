package commands

import (
	"context"
	"errors"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/libs"
	"github.com/michimani/cfkvs/types"
)

type KvsCmd struct {
	List   ListKvsSubCmd   `cmd:"" help:"List key value stores in your account."`
	Create CreateSubCmd    `cmd:"" help:"Create a key value store."`
	Delete DeleteKVSSubCmd `cmd:"" help:"Delete a key value store."`
	Info   InfoSubCmd      `cmd:"" help:"Show information of the key value store."`
	Sync   SyncSubCmd      `cmd:"" help:"Sync items in the key value store with S3 object or specified JSON file."`
}

type ListKvsSubCmd struct{}

type CreateSubCmd struct {
	Name      string `name:"name" help:"Name of the key value store." required:""`
	Comment   string `name:"comment" help:"Comment of the key value store."`
	Bucket    string `name:"bucket" help:"S3 bucket name to import key value store, if you want."`
	ObjectKey string `name:"object-key" help:"S3 object key to import key value store, if you want."`
}

type DeleteKVSSubCmd struct {
	Name string `name:"name" help:"Name of the key value store." required:""`
}

type InfoSubCmd struct {
	Name string `name:"name" help:"Name of the key value store." required:""`
}

type SyncSubCmd struct {
	Name      string `name:"name" help:"Name of the key value store." required:""`
	Bucket    string `name:"bucket" help:"S3 bucket name to sync key value store. If you want to sync with S3 object, this is required."`
	ObjectKey string `name:"object-key" help:"S3 object key to sync key value store. If you want to sync with S3 object, this is required."`
	File      string `name:"file" help:"Path to the file to sync key value store. If this is specified, sync with this file instead of S3 object."`
	Delete    bool   `name:"delete" help:"Delete items that are not in the S3 object."`
	Yes       bool   `name:"yes" short:"y" help:"Execute sync. If not specified, only show the items to be synced."`
}

func (c *ListKvsSubCmd) Run(globals *Globals) error {
	out, err := libs.ListKvs(context.TODO(), globals.CloudFrontClient)
	if err != nil {
		return err
	}

	kvsList := types.KVSList{}
	if err := kvsList.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvsList, globals.Output, globals.OutputTarget); err != nil {
		return err
	}

	return nil
}

func (c *CreateSubCmd) Run(globals *Globals) error {
	srcS3 := new(libs.KVSImportSourceS3)
	if c.Bucket != "" {
		if c.ObjectKey == "" {
			return errors.New("object-key is required when bucket is specified")
		}

		srcS3 = &libs.KVSImportSourceS3{
			Bucket: c.Bucket,
			Key:    c.ObjectKey,
		}
	}

	out, err := libs.CreateKvs(context.TODO(), globals.CloudFrontClient, c.Name, c.Comment, srcS3)
	if err != nil {
		return err
	}

	kvs := types.KVS{}
	if err := kvs.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvs, globals.Output, globals.OutputTarget); err != nil {
		return err
	}

	return nil
}

func (c *InfoSubCmd) Run(globals *Globals) error {
	info, err := libs.DescribeKeyValueStore(
		context.TODO(),
		globals.CloudFrontClient,
		globals.CloudFrontKeyValueStoreClient,
		c.Name)
	if err != nil {
		return err
	}

	if err := output.Render(info, globals.Output, globals.OutputTarget); err != nil {
		return err
	}

	return nil
}

func (c *DeleteKVSSubCmd) Run(globals *Globals) error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	if err := libs.DeleteKeyValueStore(context.TODO(), globals.CloudFrontClient, c.Name); err != nil {
		return err
	}

	_, _ = globals.OutputTarget.Write([]byte("Deleted\n"))

	return nil
}

func (c *SyncSubCmd) Run(globals *Globals) error {
	if c.Name == "" {
		return errors.New("kvs-name is required")
	}

	fromFile := false
	if c.File != "" {
		fromFile = true
	} else {
		if c.Bucket == "" {
			return errors.New("bucket is required")
		}
		if c.ObjectKey == "" {
			return errors.New("object-key is required")
		}
	}

	ctx := context.TODO()
	kvsARN, err := getKvsArn(ctx, globals.CloudFrontClient, c.Name)
	if err != nil {
		return err
	}

	// get before items
	before, err := libs.ListItems(ctx, globals.CloudFrontKeyValueStoreClient, kvsARN)
	if err != nil {
		return err
	}

	// get after items
	afterItems := &types.KeyValueStoreData{}
	if fromFile {
		// from the file
		if afterItems, err = libs.GetKeyValueStoreDataFromFile(c.File); err != nil {
			return err
		}
	} else {
		// from S3
		if afterItems, err = libs.GetKeyValueStoreData(ctx, globals.S3Client, c.Bucket, c.ObjectKey); err != nil {
			return err
		}
	}

	// show diff
	diff := before.Diff(afterItems.ToItemList(), c.Delete)
	if err := output.Render(diff, output.OutputTypeTable, globals.OutputTarget); err != nil {
		return err
	}

	if !c.Yes {
		return nil
	}

	// sync
	out, err := libs.SyncItems(ctx, globals.CloudFrontKeyValueStoreClient, kvsARN, diff.PutList(), diff.DeleteList())
	if err != nil {
		return err
	}

	kvsSimple := types.KVSSimple{}
	if err := kvsSimple.Parse(out); err != nil {
		return err
	}

	return output.Render(&kvsSimple, globals.Output, globals.OutputTarget)
}
