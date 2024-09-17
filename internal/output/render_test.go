package output_test

import (
	"bytes"
	"testing"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/types"
	"github.com/stretchr/testify/assert"
)

func Test_Render(t *testing.T) {
	cases := []struct {
		name       string
		data       any
		outputType output.OutputType
		expect     string
		wantErr    bool
	}{
		{
			name:       "ok: types.KVS, OutputTypeJson",
			data:       &types.KVS{Id: "id", Name: "name", Comment: "comment", Status: "status", ARN: "arn"},
			outputType: output.OutputTypeJson,
			expect: `{
    "id": "id",
    "name": "name",
    "comment": "comment",
    "status": "status",
    "arn": "arn"
}
`,
		},
		{
			name:       "ok: types.KVS, OutputTypeTable",
			data:       &types.KVS{Id: "id", Name: "name", Comment: "comment", Status: "status", ARN: "arn"},
			outputType: output.OutputTypeTable,
			expect: `+----+------+---------+--------+-----+
| ID | NAME | COMMENT | STATUS | ARN |
+----+------+---------+--------+-----+
| id | name | comment | status | arn |
+----+------+---------+--------+-----+
`,
		},
		{
			name:       "ok: types.Item, invalid OutputType (render as table)",
			data:       &types.Item{Key: "key", Value: "value"},
			outputType: output.OutputType("invalid"),
			expect: `+-----+-------+
| KEY | VALUE |
+-----+-------+
| key | value |
+-----+-------+
`,
		},
	}

	out := new(bytes.Buffer)

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)
			tt.Cleanup(func() {
				out.Truncate(0)
			})

			err := output.Render(c.data, c.outputType, out)
			if c.wantErr {
				asst.Error(err)
				asst.Empty(out.String())
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, out.String())
		})
	}
}
