package manifest

import (
	"cmp"
	"errors"
	"slices"

	"github.com/golang-mods/exerrors"
	"github.com/golang-mods/serrors"
	"github.com/golang-mods/sorted"
)

type Namer interface {
	Name() string
}

func ComapreName[T, U Namer](self T, other U) int {
	return cmp.Compare(self.Name(), other.Name())
}

var ErrDuplicatePath = errors.New("duplicate path")

func SortNamers[S ~[]E, E Namer](namers S) error {
	slices.SortFunc(namers, ComapreName)

	if indexes := sorted.DuplicateIndexesWith(namers, ComapreName); len(indexes) > 0 {
		return exerrors.Map0(indexes, func(index int, _ int) error {
			return serrors.Format("%w", ErrDuplicatePath)("path", namers[index].Name())
		})
	}

	return nil
}
