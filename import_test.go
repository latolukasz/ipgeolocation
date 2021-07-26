package ipgeolocation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	arg := &ImportArguments{
		DbDirectory:    "./db/",
		MysqlURI:       "root:root@tcp(localhost:3315)/ipgeolocation",
		wrongCountryID: 253, // TODO remove
	}
	err := Import(context.Background(), arg)
	assert.NoError(t, err)
}
