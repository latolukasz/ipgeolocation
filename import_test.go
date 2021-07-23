package ipgeolocation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCachedSearchLocal(t *testing.T) {
	err := Import(context.Background())
	assert.NoError(t, err)
}
