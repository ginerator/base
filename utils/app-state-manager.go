package utils

import (
	"github.com/rs/zerolog/log"
)

type Closable interface {
	Close()
}

type Monitorable interface {
	IsConnected() (bool, error)
}

type AppStateManager struct {
	closableDependencies    map[string]Closable
	monitorableDependencies map[string]Monitorable
}

func NewAppStateManager() *AppStateManager {
	return &AppStateManager{
		closableDependencies:    make(map[string]Closable),
		monitorableDependencies: make(map[string]Monitorable),
	}
}

func (manager *AppStateManager) AddClosableDependency(name string, dependency Closable) {
	manager.closableDependencies[name] = dependency
}

func (manager *AppStateManager) AddMonitorableDependency(name string, dependency Monitorable) {
	manager.monitorableDependencies[name] = dependency
}

func (manager *AppStateManager) Shutdown() {
	log.Info().Msg("Closing dependencies...")
	for _, dependency := range manager.closableDependencies {
		dependency.Close()
	}
	log.Info().Msg("All dependencies closed")
}

func (manager *AppStateManager) DependenciesConnected() (bool, error) {
	for _, dependency := range manager.monitorableDependencies {
		isConnected, err := dependency.IsConnected()
		if err != nil {
			return false, err
		}

		if !isConnected {
			return false, nil
		}
	}

	return true, nil
}
