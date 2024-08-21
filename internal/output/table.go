package output

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/jedib0t/go-pretty/table"
)

type Table struct {
	Header table.Row
	Rows   []table.Row
}

func toTable(data any) (Table, error) {
	tableData := Table{}

	switch data := data.(type) {
	case []types.KeyValueStore:
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
	}

	return tableData, nil
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
