package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createPipelineCmd = &cobra.Command{
	Use:   "create-pipeline",
	Short: "Create a new pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		jobIDs, _ := cmd.Flags().GetStringSlice("jobs")

		// Ici, vous devriez appeler l'API pour créer le pipeline
		// Pour l'exemple, nous allons simplement afficher les informations
		fmt.Printf("Creating pipeline '%s' with jobs: %v\n", name, jobIDs)
	},
}

func init() {
	rootCmd.AddCommand(createPipelineCmd)
	createPipelineCmd.Flags().String("name", "", "Name of the pipeline")
	createPipelineCmd.Flags().StringSlice("jobs", []string{}, "List of job IDs to include in the pipeline")
	createPipelineCmd.MarkFlagRequired("name")
	createPipelineCmd.MarkFlagRequired("jobs")
}
