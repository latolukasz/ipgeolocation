package ipgeolocation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	err := InitDB("./db/")
	assert.NoError(t, err)

	found, err := Search("109.173.165.107")
	assert.NoError(t, err)
	fmt.Printf("%v\n", found)
}
