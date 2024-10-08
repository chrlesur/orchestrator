package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runPipelineCmd = &cobra.Command{
	Use:   "run-pipeline",
	Short: "Run a specific pipeline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID := args[0]
		// Ici, vous devriez appeler l'API pour exécuter le pipeline
		// Pour l'exemple, nous allons simplement afficher un message
		fmt.Printf("Running pipeline: %s\n", pipelineID)
	},
}

func init() {
	rootCmd.AddCommand(runPipelineCmd)
}
