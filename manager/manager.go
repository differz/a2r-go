// Package manager ...
package manager

import (
	"fmt"
	"strings"

	"../config"
	"../data"
	"../data/cdr"
	"../data/cid"
)

// InputMethod ...
type InputMethod func() (data.In, error)

// Manager ...
type Manager struct {
	Factories map[string]InputMethod
}

// New ...
func New() Manager {
	manager := Manager{
		Factories: make(map[string]InputMethod),
	}
	manager.Register("CDR", cdr.Get)
	manager.Register("CID", cid.Get)
	return manager
}

// Register ...
func (m Manager) Register(name string, factory InputMethod) {
	_, registered := m.Factories[name]
	if !registered {
		m.Factories[name] = factory
	}
}

// Create ...
func (m Manager) Create() (data.In, error) {

	cfg := config.New()

	engineName := cfg.InputMethod
	//engineName.Valid()

	engineFactory, ok := m.Factories[engineName]
	if !ok {
		// Factory does not exist
		available := make([]string, len(m.Factories))
		for k := range m.Factories { ///+?
			available = append(available, k)
		}
		return nil, fmt.Errorf("Invalid Datastore name. Must be one of: %s", strings.Join(available, ", "))
	}

	return engineFactory()
}
