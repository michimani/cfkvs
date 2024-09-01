package output

import (
	"encoding/json"
	"os"
)

const jsonOutputIndent = "    " // 4 spaces

func RenderAsJson(data any) error {
	out, err := json.MarshalIndent(&data, "", jsonOutputIndent)
	if err != nil {
		return err
	}

	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n"))

	return nil
}
