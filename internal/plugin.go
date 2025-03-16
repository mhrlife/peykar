package internal

import (
	"errors"
	"github.com/mhrlife/peykar/pkg/common"
	"log/slog"
	"plugin"
	"strings"
)

var (
	ErrInvalidPluginDefinition = errors.New("invalid plugin definition")
)

type PluginManager struct {
	loadedPlugins map[string]common.Plugin
	logger        *slog.Logger
}

func NewPluginManager(
	logger *slog.Logger,
) *PluginManager {
	return &PluginManager{
		logger:        logger,
		loadedPlugins: make(map[string]common.Plugin),
	}
}

func (pm *PluginManager) LoadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		pm.logger.Error("Failed to load plugin", "path", path, "err", err)

		return err
	}

	newPluginSym, err := p.Lookup("NewPlugin")
	if err != nil {
		pm.logger.Error("Failed to load plugin", "path", path, "err", err)

		return err
	}

	newPlugin, ok := newPluginSym.(func() common.Plugin)
	if !ok {
		pm.logger.Error("Failed to load plugin", "path", path, "err", ErrInvalidPluginDefinition)

		return ErrInvalidPluginDefinition
	}

	instance := newPlugin()

	// extract  name from path

	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]

	pm.loadedPlugins[name] = instance

	pm.logger.Info("Loaded plugin", "path", path, "name", name)

	return nil
}

func (pm *PluginManager) All() map[string]common.Plugin {
	return pm.loadedPlugins
}
