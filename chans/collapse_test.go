package chans

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollapse(t *testing.T) {
	t.Run("should only return the most recent addition", func(t *testing.T) {
		in, out := Collapse[int](context.Background())
		defer close(in)
		in <- 1
		in <- 2
		time.Sleep(time.Second)
		assert.Equal(t, 2, <-out, "did not collapse values to most recent one")
	})

	t.Run("closing the in channel should result in the out channel also being closed", func(t *testing.T) {
		in, out := Collapse[int](context.Background())
		in <- 1
		time.Sleep(time.Second)
		close(in)
		time.Sleep(time.Second)
		_, open := <- out
		assert.False(t, open, "channel was still open")
	})
}
