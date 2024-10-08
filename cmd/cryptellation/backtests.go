package main

import (
	"github.com/spf13/cobra"
)

var backtestsCmd = &cobra.Command{
	Use:     "backtests",
	Aliases: []string{"bt"},
	Short:   "Manipulate backtests service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := executeParentPersistentPreRuns(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

var backtestsInfoCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Read info from backtests service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return displayServiceInfo(globalClient.Backtests())
	},
}

func initBacktests(rootCmd *cobra.Command) {
	backtestsCmd.AddCommand(backtestsInfoCmd)
	rootCmd.AddCommand(backtestsCmd)
}
