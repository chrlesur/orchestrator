package commands

import (
	"fmt"

	"orchestrator/internal/storage"

	"github.com/spf13/cobra"
)

var listPipelinesCmd = &cobra.Command{
	Use:   "list-pipelines",
	Short: "List all pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		pipelines, err := storage.ListPipelines()
		if err != nil {
			fmt.Printf("Error listing pipelines: %v\n", err)
			return
		}

		fmt.Println("Pipelines:")
		for _, pipeline := range pipelines {
			fmt.Printf("- ID: %s, Name: %s, Status: %s\n", pipeline.ID, pipeline.Name, pipeline.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(listPipelinesCmd)
}
