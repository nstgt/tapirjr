package main

import (
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "tapirjr",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	receiverCmd := newReceiverCmd()
	senderCmd := newSenderCmd()
	rootCmd.AddCommand(receiverCmd, senderCmd)

	return rootCmd
}
