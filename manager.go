package tango

import (
	"log/slog"
	"path/filepath"

	"github.com/golang-mods/tango/internal/gocommand"
	"github.com/spf13/afero"
)

var (
	fs      = afero.NewOsFs()
	command = gocommand.NewMemorized(gocommand.NewEntity())
)

type Manager struct {
	current string
	logger  *slog.Logger
}

func NewManager(current string, logger *slog.Logger) (*Manager, error) {
	abs, err := filepath.Abs(current)
	if err != nil {
		return nil, err
	}

	return &Manager{current: abs, logger: logger}, nil
}
