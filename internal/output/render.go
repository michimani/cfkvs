package output

import "os"

type OutputType string

const (
	OutputTypeJson  OutputType = "json"
	OutputTypeTable OutputType = "table"
)

func Render(data any, outputType OutputType) error {
	switch outputType {
	case OutputTypeJson:
		return RenderAsJson(data, os.Stdout)
	case OutputTypeTable:
		return RenderAsTable(data, os.Stdout)
	default:
		return RenderAsTable(data, os.Stdout)
	}
}
