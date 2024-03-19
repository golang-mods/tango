package memorize

import "sync"

func Memorized[I comparable, O any](callback func(I) O) func(I) O {
	var memo sync.Map

	return func(input I) O {
		if output, ok := memo.Load(input); ok {
			return output.(O)
		}

		output := callback(input)
		memo.Store(input, output)
		return output
	}
}
