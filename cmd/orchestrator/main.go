package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/chrlesur/orchestrator/internal/api"
	"github.com/chrlesur/orchestrator/internal/config"
	"github.com/chrlesur/orchestrator/internal/db"
	"github.com/chrlesur/orchestrator/internal/job"
	"github.com/chrlesur/orchestrator/internal/pipeline"
	"github.com/chrlesur/orchestrator/internal/plugin"
	"github.com/chrlesur/orchestrator/internal/ui"
	"github.com/chrlesur/orchestrator/pkg/logger"
	"github.com/chrlesur/orchestrator/pkg/version"
)

func main() {
	// Afficher la version
	fmt.Printf("Orchestrator version %s\n", version.GetVersion())

	// Charger la configuration
	cfg, err := config.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la configuration: %v", err)
	}

	// Initialiser le logger
	err = logger.Init(cfg.Logging.Level, cfg.Logging.File)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du logger: %v", err)
	}

	// Initialiser la base de données BoltDB
	store, err := db.NewStore(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données: %v", err)
	}
	defer store.Close()

	// Initialiser le gestionnaire de plugins
	pluginManager := plugin.NewPluginManager()

	// Charger les plugins
	pluginsDir := "./plugins"
	err = filepath.Walk(pluginsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".so" {
			if err := pluginManager.LoadPlugin(path); err != nil {
				logger.Error(fmt.Sprintf("Failed to load plugin %s: %v", path, err))
			}
		}
		return nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error walking the plugins directory: %v", err))
	}

	// Créer le gestionnaire de jobs
	jobManager := job.NewManager(5, store, cfg.Jobs.DefaultTimeout, cfg.Jobs.MaxRetries, pluginManager)

	// Créer le gestionnaire de pipelines
	pipelineManager := pipeline.NewManager(3, store, pluginManager)

	// Créer et lancer l'interface TUI dans une goroutine
	tui := ui.NewTUI(jobManager, pipelineManager, pluginManager)
	go func() {
		if err := tui.Run(); err != nil {
			log.Fatalf("Erreur lors de l'exécution de l'interface TUI: %v", err)
		}
	}()

	// Créer et lancer le serveur API dans une goroutine
	apiServer := api.NewServer(jobManager, pipelineManager, pluginManager)
	go func() {
		logger.Info(fmt.Sprintf("Démarrage du serveur API sur :%d", cfg.Server.Port))
		if err := apiServer.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			log.Fatalf("Erreur lors du démarrage du serveur API: %v", err)
		}
	}()

	// Attendre un signal d'arrêt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Arrêter proprement les gestionnaires
	jobManager.Shutdown()
	pipelineManager.Shutdown()

	logger.Info("Application terminée")
}
