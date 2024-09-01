package commands

import (
	"context"
	"errors"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/libs"
	"github.com/michimani/cfkvs/types"
)

type KvsCmd struct {
	List   ListKvsSubCmd `cmd:"" help:"List key value stores in your account."`
	Create CreateSubCmd  `cmd:"" help:"Create a key value store."`
}

type ListKvsSubCmd struct{}

type CreateSubCmd struct {
	Name      string `name:"name" help:"Name of the key value store." required:""`
	Comment   string `name:"comment" help:"Comment of the key value store."`
	Bucket    string `name:"bucket" help:"S3 bucket name to import key value store, if you want."`
	ObjectKey string `name:"object-key" help:"S3 object key to import key value store, if you want."`
}

func (c *ListKvsSubCmd) Run(globals *Globals) error {
	cfc, err := libs.NewCloudFrontClient(context.TODO())
	if err != nil {
		return err
	}

	out, err := libs.ListKvs(context.TODO(), cfc)
	if err != nil {
		return err
	}

	kvsList := types.KVSList{}
	if err := kvsList.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvsList, globals.Output); err != nil {
		return err
	}

	return nil
}

func (c *CreateSubCmd) Run(globals *Globals) error {
	cfc, err := libs.NewCloudFrontClient(context.TODO())
	if err != nil {
		return err
	}

	var srcS3 *libs.KVSImportSourceS3 = nil
	if c.Bucket != "" {
		if c.ObjectKey == "" {
			return errors.New("object-key is required when bucket is specified")
		}

		srcS3 = &libs.KVSImportSourceS3{
			Bucket: c.Bucket,
			Key:    c.ObjectKey,
		}
	}

	out, err := libs.CreateKvs(context.TODO(), cfc, c.Name, c.Comment, srcS3)
	if err != nil {
		return err
	}

	kvs := types.KVS{}
	if err := kvs.Parse(out); err != nil {
		return err
	}

	if err := output.Render(&kvs, globals.Output); err != nil {
		return err
	}

	return nil
}
