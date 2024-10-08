package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getContextCmd = &cobra.Command{
	Use:   "get-context [pipeline|job] [id]",
	Short: "Get context of a pipeline or job",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		contextType := args[0]
		id := args[1]

		switch contextType {
		case "pipeline":
			fmt.Printf("Getting context for pipeline: %s\n", id)
			// Appeler l'API pour récupérer le contexte du pipeline
		case "job":
			fmt.Printf("Getting context for job: %s\n", id)
			// Appeler l'API pour récupérer le contexte du job
		default:
			fmt.Println("Invalid context type. Use 'pipeline' or 'job'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(getContextCmd)
}
