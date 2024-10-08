package commands

import (
    "fmt"

    "github.com/spf13/cobra"
    "orchestrator/internal/logging"
	"orchestrator/internal/storage"
)

var runJobCmd = &cobra.Command{
	Use:   "run-job",
	Short: "Run a specific job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobID := args[0]
		
		job, err := storage.GetJob(jobID)
		if err != nil {
			fmt.Printf("Error retrieving job: %v\n", err)
			return
		}

		fmt.Printf("Running job: %s\n", job.Name)
		// Ici, vous devriez implémenter la logique pour exécuter réellement le job
		// Pour l'exemple, nous allons simplement mettre à jour son statut
		job.Status = "RUNNING"
		err = storage.SaveJob(job)
		if err != nil {
			fmt.Printf("Error updating job status: %v\n", err)
			return
		}

		fmt.Println("Job execution started.")
	},
}

func init() {
	rootCmd.AddCommand(runJobCmd)
}

func runRunJob(cmd *cobra.Command, args []string) {
    err := apiClient.RunJob(args[0])
    if err != nil {
        logging.ErrorLogger.Printf("Failed to run job: %v", err)
        fmt.Printf("Failed to run job: %v\n", err)
        return
    }

    fmt.Println("Job started successfully")
}