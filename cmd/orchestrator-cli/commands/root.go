package commands

import (
	"fmt"
	"os"

	"orchestrator/internal/config"
	"orchestrator/internal/constants"
	"orchestrator/internal/logging"

	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	silent    bool
	debug     bool
	apiClient *api.Client
)

var rootCmd = &cobra.Command{
	Use:   "orchestrator-cli",
	Short: "ORCHESTRATOR CLI client",
	Long:  `Command-line interface for interacting with the ORCHESTRATOR server.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize config, logging, and API client
		initializeClient()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "run in silent mode (no console output)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "run in debug mode (verbose logging)")
}

func initializeClient() {
	// Load configuration
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	initializeLogger(cfg)

	// Initialize API client
	apiClient = api.NewClient(fmt.Sprintf("http://localhost:%d", cfg.Server.Port))
}

func initializeLogger(cfg *config.Config) {
	if silent {
		cfg.Logging.ToConsole = false
	}
	if debug {
		cfg.Logging.DebugMode = true
	}

	logFile := fmt.Sprintf("logs/cli-%s", cfg.Logging.File)
	err := logging.InitLoggers(logFile, cfg.Logging.ToConsole, cfg.Logging.DebugMode)
	if err != nil {
		fmt.Printf("Failed to initialize loggers: %v\n", err)
		os.Exit(1)
	}

	logging.InfoLogger.Printf("ORCHESTRATOR CLI v%s initialized", constants.Version)
	logging.InfoLogger.Printf("Logging to console: %v", cfg.Logging.ToConsole)
	logging.InfoLogger.Printf("Debug mode: %v", cfg.Logging.DebugMode)
}
