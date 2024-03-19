package tango

import (
	"os"
	"path/filepath"

	"github.com/golang-mods/tango/internal/constant"
)

func (manager *Manager) Init() error {
	name := filepath.Join(manager.current, constant.ManifestFileName)

	file, err := fs.OpenFile(name, os.O_CREATE|os.O_EXCL, 0600)
	if os.IsExist(err) {
		manager.logger.Warn("file already exists", "name", name)
		return nil
	} else if err != nil {
		return err
	}
	defer file.Close()

	manager.logger.Info("create file", "name", name)

	return nil
}
