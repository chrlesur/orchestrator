package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getPipelineCmd = &cobra.Command{
	Use:   "get-pipeline",
	Short: "Get details of a specific pipeline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID := args[0]
		// Ici, vous devriez appeler l'API pour obtenir les détails du pipeline
		// Pour l'exemple, nous allons simplement afficher l'ID
		fmt.Printf("Getting details for pipeline: %s\n", pipelineID)
	},
}

func init() {
	rootCmd.AddCommand(getPipelineCmd)
}
