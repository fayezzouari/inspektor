package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"inspektor/internal/inspector"
)

var rootCmd = &cobra.Command{
	Use:   "inspektor [PID]",
	Short: "AI-powered process inspector and system monitor",
	Long: `Inspektor analyzes running processes and system resources,
providing detailed insights and AI-generated warnings about system health.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid PID: %s\n", args[0])
			os.Exit(1)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		verbose, _ := cmd.Flags().GetBool("verbose")

		inspector := inspector.New()
		if err := inspector.InspectWithOptions(int32(pid), jsonOutput, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error inspecting process: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}