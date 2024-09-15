package output

import (
	"encoding/json"
	"io"
)

const jsonOutputIndent = "    " // 4 spaces

func RenderAsJson(data any, o io.Writer) error {
	out, err := json.MarshalIndent(&data, "", jsonOutputIndent)
	if err != nil {
		return err
	}

	// In most cases, the Write method does not return an error, so error checking is omitted.
	_, _ = o.Write(out)
	_, _ = o.Write([]byte("\n"))

	return nil
}
