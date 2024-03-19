package manifest

import (
	"bytes"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestManifest(t *testing.T) {
	toml := `[tools]
'example.com/a' = 'v1.1'
'example.com/b/c' = {version = 'v2', env = ['foo=bar']}
`

	manifest := Manifest{
		Tools: []ManifestTool{
			{
				ToolPath: ToolPath{Path: "example.com/a"},
				Version:  *lo.Must(NewConstraints("v1.1")),
			},
			{
				ToolPath: ToolPath{Path: "example.com/b/c"},
				Version:  *lo.Must(NewConstraints("v2")),
				Tool:     Tool{Env: []string{"foo=bar"}},
			},
		},
	}

	t.Run("EncodeManifest", func(t *testing.T) {
		var buffer bytes.Buffer

		assert.Nil(t, EncodeManifest(&buffer, &manifest))
		assert.Equal(t, toml, buffer.String())
	})

	t.Run("DecodeManifest", func(t *testing.T) {
		actual, err := DecodeManifest(bytes.NewBufferString(toml))

		assert.Nil(t, err)
		assert.Equal(t, &manifest, actual)
	})
}
