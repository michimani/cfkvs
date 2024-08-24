package output

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	kvs "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvsTypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
	"github.com/jedib0t/go-pretty/table"
)

type Table struct {
	Header table.Row
	Rows   []table.Row
}

func toTable(data any) (*Table, error) {
	tableData := Table{}

	switch data := data.(type) {
	case []cfTypes.KeyValueStore:
		// List of CloudFront Key Value Stores
		tableData.Header = table.Row{"ID", "Name", "Comment", "Status", "ARN"}
		for _, kvs := range data {
			tableData.Rows = append(
				tableData.Rows,
				table.Row{
					aws.ToString(kvs.Id),
					aws.ToString(kvs.Name),
					aws.ToString(kvs.Comment),
					aws.ToString(kvs.Status),
					aws.ToString(kvs.ARN)})
		}

	case []kvsTypes.ListKeysResponseListItem:
		// List of Items in the Key Value Store
		tableData.Header = table.Row{"Key", "Value"}
		for _, item := range data {
			tableData.Rows = append(
				tableData.Rows,
				table.Row{
					aws.ToString(item.Key),
					aws.ToString(item.Value)})
		}

	case *kvs.GetKeyOutput:
		// An item in the Key Value Store
		tableData.Header = table.Row{"Key", "Value"}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{
				aws.ToString(data.Key),
				aws.ToString(data.Value)})

	case *kvs.PutKeyOutput:
		// Put an item in the Key Value Store
		tableData.Header = table.Row{"ItemCount", "TotalSize (bytes)"}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{
				aws.ToInt32(data.ItemCount),
				aws.ToInt64(data.TotalSizeInBytes)})

	case *kvs.DeleteKeyOutput:
		// Delete an item in the Key Value Store
		tableData.Header = table.Row{"ItemCount", "TotalSize (bytes)"}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{
				aws.ToInt32(data.ItemCount),
				aws.ToInt64(data.TotalSizeInBytes)})

	default:
		return nil, errors.New("unsupported data type")
	}

	return &tableData, nil
}

func RenderAsTable(data any) error {
	tableData, err := toTable(data)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(tableData.Header)
	t.AppendRows(tableData.Rows)
	t.Render()

	return nil
}
