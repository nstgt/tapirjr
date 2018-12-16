package main

import (
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "tapirjr",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	serverCmd := newServerCmd()
	clientCmd := newClientCmd()
	rootCmd.AddCommand(serverCmd, clientCmd)

	return rootCmd
}
