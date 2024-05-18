package chans

import (
	"context"
)

// Collapse will take any number of objects on T, but reading from the output will only
// return the most recent one, it is the writer's responsibility to close the `in` channel
func Collapse[T any](ctx context.Context) (in chan<- T, out <-chan T) {
	inChan, outChan := make(chan T), make(chan T)
	go func() {
		defer close(outChan)
		for {
			select {
			case <-ctx.Done():
				return
			case thing, stillOpen := <-inChan:
				if !stillOpen {
					return
				}
			inner:
				for {
					select {
					case <-ctx.Done():
						return
					case outChan <- thing:
						break inner
					case nextThing, stillOpen := <-inChan:
						if !stillOpen {
							return
						}
						thing = nextThing
					}
				}
			}
		}
	}()
	return inChan, outChan
}
