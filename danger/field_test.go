package danger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetPrivateField(t *testing.T) {
	type tstStruct struct{ y int }
	x := tstStruct{y: 1}
	SetPrivateField(&x, "y", 4)
	assert.Equal(t, x.y, 4)
}
