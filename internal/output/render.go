package output

import "io"

type OutputType string

const (
	OutputTypeJson  OutputType = "json"
	OutputTypeTable OutputType = "table"
)

func Render(data any, outputType OutputType, w io.Writer) error {
	switch outputType {
	case OutputTypeJson:
		return RenderAsJson(data, w)
	case OutputTypeTable:
		return RenderAsTable(data, w)
	default:
		return RenderAsTable(data, w)
	}
}
