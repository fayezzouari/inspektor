package cmd

import (
	"fmt"
	"os"
	"strconv"

	"inspektor/internal/inspector"

	"github.com/spf13/cobra"
)

var (
	portFlag int
)

var rootCmd = &cobra.Command{
	Use:   "inspektor [PID]",
	Short: "AI-powered process inspector and system monitor",
	Long: `Inspektor analyzes running processes and system resources,
providing detailed insights and AI-generated warnings about system health.

You can inspect a process by:
  - PID: inspektor 1234
  - Port: inspektor --port 8080`,
	Args: func(cmd *cobra.Command, args []string) error {
		// If port flag is set, no args needed
		if portFlag > 0 {
			return nil
		}
		// Otherwise, require exactly one PID argument
		if len(args) != 1 {
			return fmt.Errorf("requires either a PID argument or --port flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		verbose, _ := cmd.Flags().GetBool("verbose")

		insp := inspector.New()

		var err error
		if portFlag > 0 {
			// Inspect by port
			err = insp.InspectByPort(portFlag, jsonOutput, verbose)
		} else {
			// Inspect by PID
			pid, parseErr := strconv.Atoi(args[0])
			if parseErr != nil {
				fmt.Fprintf(os.Stderr, "Invalid PID: %s\n", args[0])
				os.Exit(1)
			}
			err = insp.InspectWithOptions(int32(pid), jsonOutput, verbose)
		}

		if err != nil {
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
	rootCmd.Flags().IntVarP(&portFlag, "port", "p", 0, "Inspect process listening on specified port")
}
