package commands

type KvsCmd struct {
	List   ListKvsSubCmd `cmd:"" help:"List key value stores in your account."`
	Create CreateSubCmd  `cmd:"" help:"Create a key value store."`
}

type ListKvsSubCmd struct{}

type CreateSubCmd struct {
	Name    string `arg:"" name:"name" help:"Name of the key value store." required:""`
	Comment string `arg:"" name:"comment" help:"Comment of the key value store."`
}
