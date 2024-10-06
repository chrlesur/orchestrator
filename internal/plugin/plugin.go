package plugin

import (
	"fmt"
	"plugin"

	"github.com/chrlesur/orchestrator/pkg/logger"
)

// Plugin définit l'interface que tous les plugins doivent implémenter
type Plugin interface {
	Name() string
	Version() string
	Execute(args map[string]interface{}) (interface{}, error)
}

// PluginManager gère le chargement et l'exécution des plugins
type PluginManager struct {
	plugins map[string]Plugin
}

// NewPluginManager crée une nouvelle instance de PluginManager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

// LoadPlugin charge un plugin à partir d'un fichier .so
func (pm *PluginManager) LoadPlugin(path string) error {
	// Ouvrir le plugin
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("could not open plugin: %v", err)
	}

	// Chercher le symbole "New"
	newFunc, err := p.Lookup("New")
	if err != nil {
		return fmt.Errorf("could not find 'New' symbol: %v", err)
	}

	// Vérifier que le symbole est une fonction qui retourne un Plugin
	newPlugin, ok := newFunc.(func() Plugin)
	if !ok {
		return fmt.Errorf("'New' symbol is not a function that returns a Plugin")
	}

	// Créer une instance du plugin
	instance := newPlugin()

	// Enregistrer le plugin
	pm.plugins[instance.Name()] = instance

	logger.Info(fmt.Sprintf("Loaded plugin: %s (version %s)", instance.Name(), instance.Version()))

	return nil
}

// ExecutePlugin exécute un plugin chargé
func (pm *PluginManager) ExecutePlugin(name string, args map[string]interface{}) (interface{}, error) {
	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin.Execute(args)
}

// GetLoadedPlugins retourne la liste des plugins chargés
func (pm *PluginManager) GetLoadedPlugins() []string {
	var pluginNames []string
	for name := range pm.plugins {
		pluginNames = append(pluginNames, name)
	}
	return pluginNames
}
