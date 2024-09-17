package output_test

import (
	"bytes"
	"testing"

	"github.com/michimani/cfkvs/internal/output"
	"github.com/michimani/cfkvs/types"
	"github.com/stretchr/testify/assert"
)

func Test_RenderAsJson(t *testing.T) {
	cases := []struct {
		name    string
		data    any
		expect  string
		wantErr bool
	}{
		{
			name: "ok",
			data: types.KVS{
				Id:      "id",
				Name:    "name",
				Comment: "comment",
				Status:  "status",
				ARN:     "arn",
			},
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
			name:    "invalid json",
			data:    make(chan int),
			wantErr: true,
		},
	}

	out := new(bytes.Buffer)

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			tt.Cleanup(func() {
				out.Truncate(0)
			})

			asst := assert.New(tt)

			err := output.RenderAsJson(c.data, out)
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
