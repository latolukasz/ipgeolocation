package ipgeolocation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	arg := &ImportArguments{
		DbDirectory: "./db/",
	}
	err := Import(arg)
	assert.NoError(t, err)
}
