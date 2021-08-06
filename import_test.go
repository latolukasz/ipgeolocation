package ipgeolocation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	err := Import("./db/")
	assert.NoError(t, err)
}
