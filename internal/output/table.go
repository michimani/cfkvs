package output

import (
	"errors"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/michimani/cfkvs/types"
)

type Table struct {
	Header table.Row
	Rows   []table.Row
}

func toTable(data any) (*Table, error) {
	tableData := Table{}

	switch data := data.(type) {
	case *types.KVS:
		// CloudFront Key Value Store
		tableData.Header = table.Row{"ID", "Name", "Comment", "Status", "ARN"}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{data.Id, data.Name, data.Comment, data.Status, data.ARN})

	case *types.KVSList:
		// List of CloudFront Key Value Stores
		tableData.Header = table.Row{"ID", "Name", "Comment", "Status", "ARN"}
		for _, kvs := range *data {
			tableData.Rows = append(
				tableData.Rows,
				table.Row{kvs.Id, kvs.Name, kvs.Comment, kvs.Status, kvs.ARN})
		}

	case *types.ItemList:
		// List of Items in the Key Value Store
		tableData.Header = table.Row{"Key", "Value"}
		for _, item := range *data {
			tableData.Rows = append(
				tableData.Rows,
				table.Row{item.Key, item.Value})
		}

	case *types.Item:
		// An item in the Key Value Store
		tableData.Header = table.Row{"Key", "Value"}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{data.Key, data.Value})

	case *types.KVSSimple:
		// Put an item in the Key Value Store
		tableData.Header = table.Row{"ItemCount", "TotalSize (bytes)"}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{data.ItemCount, data.TotalSize})

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
