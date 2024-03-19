package gocommand

import (
	"github.com/golang-mods/tango/internal/memorize"
	"github.com/samber/lo"
)

type Memorized struct {
	Command
	memorizedVersions func(string) lo.Tuple2[*Versions, error]
}

func NewMemorized(command Command) *Memorized {
	return &Memorized{
		Command: command,
		memorizedVersions: memorize.Memorized(func(pkg string) lo.Tuple2[*Versions, error] {
			return lo.T2(command.Versions(pkg))
		}),
	}
}

func (command Memorized) MemorizedVersions(pkg string) (*Versions, error) {
	return command.memorizedVersions(pkg).Unpack()
}
