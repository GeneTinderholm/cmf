package cmf

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckErr(t *testing.T) {
	t.Run("should panic if given an error", func(t *testing.T) {
		defer func() {
			recover()
		}()
		CheckErr(errors.New("should panic"))
		t.Fail()
	})
	t.Run("should not panic if given nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fail()
			}
		}()
		CheckErr(nil)
	})
}

func TestMust(t *testing.T) {
	t.Run("should convert an error tuple into a non-error value", func(t *testing.T) {
		mightFail := func() (int, error) { return 1, nil }

		i := Must(mightFail())

		assert.Equal(t, 1, i)
	})
	t.Run("should panic on an error", func(t *testing.T) {
		defer func() {
			recover()
		}()
		mightFail := func() (int, error) { return 0, errors.New("should panic") }

		_ = Must(mightFail())

		t.Fail()
	})
}
