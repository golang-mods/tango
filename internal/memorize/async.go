package memorize

import (
	"sync"

	"github.com/samber/lo"
)

func Async[T any](callback func() (T, error)) func() (T, error) {
	var result T
	var err error
	channel := lo.Async2(callback)
	once := sync.OnceFunc(func() {
		result, err = (<-channel).Unpack()
	})

	return func() (T, error) {
		once()
		return result, err
	}
}
