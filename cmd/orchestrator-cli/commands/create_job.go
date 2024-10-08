package commands

import (
	"fmt"

	"orchestrator/internal/job"
	"orchestrator/internal/logging"
	"orchestrator/internal/storage"

	"github.com/spf13/cobra"
)

var createJobCmd = &cobra.Command{
	Use:   "create-job",
	Short: "Create a new job",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		command, _ := cmd.Flags().GetString("command")

		newJob := job.NewJob(name, command)
		err := storage.SaveJob(newJob)
		if err != nil {
			fmt.Printf("Error creating job: %v\n", err)
			return
		}

		fmt.Printf("Job created successfully. ID: %s\n", newJob.ID)
	},
}

func init() {
	rootCmd.AddCommand(createJobCmd)
	createJobCmd.Flags().String("name", "", "Name of the job")
	createJobCmd.Flags().String("command", "", "Command to run")
	createJobCmd.MarkFlagRequired("name")
	createJobCmd.MarkFlagRequired("command")
}

func runCreateJob(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	command, _ := cmd.Flags().GetString("command")
	workDir, _ := cmd.Flags().GetString("work-dir")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	maxRetries, _ := cmd.Flags().GetInt("max-retries")

	job, err := apiClient.CreateJob(name, command, workDir, timeout, maxRetries)
	if err != nil {
		logging.ErrorLogger.Printf("Failed to create job: %v", err)
		fmt.Printf("Failed to create job: %v\n", err)
		return
	}

	fmt.Printf("Job created successfully: %v\n", job)
}
