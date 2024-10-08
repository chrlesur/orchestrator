package commands

import (
	"fmt"

	"orchestrator/internal/logging"
	"orchestrator/internal/storage"

	"github.com/spf13/cobra"
)

var listJobsCmd = &cobra.Command{
	Use:   "list-jobs",
	Short: "List all jobs",
	Run: func(cmd *cobra.Command, args []string) {
		jobs, err := storage.ListJobs()
		if err != nil {
			fmt.Printf("Error listing jobs: %v\n", err)
			return
		}

		fmt.Println("Jobs:")
		for _, job := range jobs {
			fmt.Printf("- ID: %s, Name: %s, Status: %s\n", job.ID, job.Name, job.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(listJobsCmd)
}

func runListJobs(cmd *cobra.Command, args []string) {
	jobs, err := apiClient.ListJobs()
	if err != nil {
		logging.ErrorLogger.Printf("Failed to list jobs: %v", err)
		fmt.Printf("Failed to list jobs: %v\n", err)
		return
	}

	for _, job := range jobs {
		fmt.Printf("Job ID: %s, Name: %s, Status: %s\n", job["ID"], job["Name"], job["Status"])
	}
}
