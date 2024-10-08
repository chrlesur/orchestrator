package main

import (
	"fmt"
	"orchestrator/cmd/orchestrator-server/server"
	"orchestrator/internal/config"
	"orchestrator/internal/constants"
	"orchestrator/internal/logging"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	silent  bool
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "orchestrator-server",
	Short: "ORCHESTRATOR server application",
	Long:  `ORCHESTRATOR server manages jobs, pipelines, and schedules.`,
	Run:   runServer,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "run in silent mode (no console output)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "run in debug mode (verbose logging)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	cfg := loadConfig()
	initializeLogger(cfg)
	srv := server.NewServer(cfg)
	srv.Run()
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

func initializeLogger(cfg *config.Config) {
	if silent {
		cfg.Logging.ToConsole = false
	}
	if debug {
		cfg.Logging.DebugMode = true
	}

	logFile := fmt.Sprintf("logs/%s", cfg.Logging.File)
	err := logging.InitLoggers(logFile, cfg.Logging.ToConsole, cfg.Logging.DebugMode)
	if err != nil {
		fmt.Printf("Failed to initialize loggers: %v\n", err)
		os.Exit(1)
	}

	logging.InfoLogger.Printf("Initializing ORCHESTRATOR v%s components...", constants.Version)
	logging.InfoLogger.Printf("Server will run on port: %d", cfg.Server.Port)
	logging.InfoLogger.Printf("Database path: %s", cfg.Database.Path)
	logging.InfoLogger.Printf("Log level: %s", cfg.Logging.Level)
	logging.InfoLogger.Printf("Logging to console: %v", cfg.Logging.ToConsole)
	logging.InfoLogger.Printf("Debug mode: %v", cfg.Logging.DebugMode)
}
