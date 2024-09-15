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

	if _, err := o.Write(out); err != nil {
		return err
	}

	_, _ = o.Write([]byte("\n"))

	return nil
}
