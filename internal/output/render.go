package output

type OutputType string

const (
	OutputTypeJson  OutputType = "json"
	OutputTypeTable OutputType = "table"
)

func Render(data any, outputType OutputType) error {
	switch outputType {
	case OutputTypeJson:
		return RenderAsJson(data)
	case OutputTypeTable:
		return RenderAsTable(data)
	default:
		return RenderAsTable(data)
	}
}
