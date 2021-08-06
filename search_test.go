package ipgeolocation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	err := InitDB("./db/")
	assert.NoError(t, err)
	found, err := Search("37.143.210.32")
	assert.NoError(t, err)
	assert.Equal(t, "Sofia", found.City)
}
