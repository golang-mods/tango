package tango

import (
	"github.com/golang-mods/tango/internal/process"
	"github.com/golang-mods/tango/internal/process/binary"
)

type Binary = binary.Binary

func (manager *Manager) List() ([]*Binary, error) {
	return process.ReadBinariesDirectory(fs, command, manager.current)
}
