package output

import (
	"fmt"
	"io"
	"reflect"

	"github.com/jedib0t/go-pretty/table"
	"github.com/michimani/cfkvs/types"
)

type Table struct {
	// Descriptions of the table
	// These descriptions will render raw text before the table
	Descriptions []string

	Headers []table.Row
	Rows    []table.Row
}

func toTables(data any) ([]Table, error) {
	tableData := Table{}

	switch data := data.(type) {
	case *types.KVS:
		// CloudFront Key Value Store
		tableData.Headers = []table.Row{{"ID", "Name", "Comment", "Status", "ARN"}}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{data.Id, data.Name, data.Comment, data.Status, data.ARN})
		return []Table{tableData}, nil

	case *types.KVSList:
		// List of CloudFront Key Value Stores
		tableData.Headers = []table.Row{{"ID", "Name", "Comment", "Status", "ARN"}}
		for _, kvs := range *data {
			tableData.Rows = append(
				tableData.Rows,
				table.Row{kvs.Id, kvs.Name, kvs.Comment, kvs.Status, kvs.ARN})
		}
		return []Table{tableData}, nil

	case *types.ItemList:
		// List of Items in the Key Value Store
		tableData.Headers = []table.Row{{"Key", "Value"}}
		for _, item := range data.Data {
			tableData.Rows = append(
				tableData.Rows,
				table.Row{item.Key, item.Value})
		}
		return []Table{tableData}, nil

	case *types.Item:
		// An item in the Key Value Store
		tableData.Headers = []table.Row{{"Key", "Value"}}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{data.Key, data.Value})
		return []Table{tableData}, nil

	case *types.KVSSimple:
		// Put an item in the Key Value Store
		tableData.Headers = []table.Row{{"ItemCount", "TotalSizeInBytes"}}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{data.ItemCount, data.TotalSize})
		return []Table{tableData}, nil

	case *types.KeyValueStoreFull:
		// Full information of the Key Value Store
		tableData.Headers = []table.Row{{
			"ID",
			"ARN",
			"Name",
			"Comment",
			"Status",
			"ItemCount",
			"TotalSizeInBytes",
			"Created",
			"LastModified",
			"FailureReason",
			"ETag",
		}}
		tableData.Rows = append(
			tableData.Rows,
			table.Row{
				data.ID,
				data.ARN,
				data.Name,
				data.Comment,
				data.Status,
				data.ItemCount,
				data.TotalSizeInBytes,
				data.Created,
				data.LastModified,
				data.FailureReason,
				data.ETag,
			})
		return []Table{tableData}, nil

	case *types.ItemListDiff:
		// List of Item differences in the Key Value Store
		tables := []Table{}

		// Add
		if len(data.Add) == 0 {
			tables = append(tables, Table{
				Descriptions: []string{"\n[ADDED] No items will be added."},
			})
		} else {
			addHeaders := []table.Row{
				{"#", "Key", "Value"},
			}
			addRows := []table.Row{}
			for i, diff := range data.Add {
				addRows = append(addRows, table.Row{i + 1, diff.After.Key, diff.After.Value})
			}
			tables = append(tables, Table{
				Descriptions: []string{"\n[ADDED] Following items will be added."},
				Headers:      addHeaders,
				Rows:         addRows})
		}

		// Update
		if len(data.Update) == 0 {
			tables = append(tables, Table{
				Descriptions: []string{"\n[UPDATED] No items will be updated."},
			})
		} else {
			updateHeaders := []table.Row{
				{"#", "Key", "Before Value", "After Value"},
			}

			updateRows := []table.Row{}
			for i, diff := range data.Update {
				updateRows = append(updateRows, table.Row{i + 1, diff.Before.Key, diff.Before.Value, diff.After.Value})
			}
			tables = append(tables, Table{
				Descriptions: []string{"\n[UPDATED] Following items will be updated."},
				Headers:      updateHeaders,
				Rows:         updateRows})
		}

		// Delete
		if len(data.Delete) == 0 {
			tables = append(tables, Table{
				Descriptions: []string{"\n[DELETED] No items will be deleted."},
			})
		} else {
			deleteHeaders := []table.Row{
				{"#", "Key", "Value"},
			}
			deleteRows := []table.Row{}
			for i, diff := range data.Delete {
				deleteRows = append(deleteRows, table.Row{i + 1, diff.Before.Key, diff.Before.Value})
			}
			tables = append(tables, Table{
				Descriptions: []string{"\n[DELETED] Following items will be deleted."},
				Headers:      deleteHeaders,
				Rows:         deleteRows})
		}

		return tables, nil

	default:
		return nil, fmt.Errorf("failed to render as table due to unexpected type: %s", reflect.TypeOf(data).String())
	}
}

func RenderAsTable(data any, o io.Writer) error {
	tables, err := toTables(data)
	if err != nil {
		return err
	}

	for _, tableData := range tables {
		// render descriptions
		for _, desc := range tableData.Descriptions {
			_, _ = o.Write([]byte(desc + "\n"))
		}

		if len(tableData.Headers) == 0 || len(tableData.Rows) == 0 {
			continue
		}

		t := table.NewWriter()
		t.SetOutputMirror(o)
		for _, header := range tableData.Headers {
			t.AppendHeader(header)
		}
		t.AppendRows(tableData.Rows)
		t.Render()
	}

	return nil
}
