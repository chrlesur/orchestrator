package commands

import (
	"fmt"

	"orchestrator/internal/logging"

	"github.com/spf13/cobra"
)

var getJobCmd = &cobra.Command{
	Use:   "get-job [id]",
	Short: "Get details of a specific job",
	Args:  cobra.ExactArgs(1),
	Run:   runGetJob,
}

func init() {
	rootCmd.AddCommand(getJobCmd)
}

func runGetJob(cmd *cobra.Command, args []string) {
	job, err := apiClient.GetJob(args[0])
	if err != nil {
		logging.ErrorLogger.Printf("Failed to get job: %v", err)
		fmt.Printf("Failed to get job: %v\n", err)
		return
	}

	fmt.Printf("Job details: %v\n", job)
}
