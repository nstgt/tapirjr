package main

import (
	"github.com/nstgt/tapirjr/internal/pkg/receiver"
	"github.com/spf13/cobra"
)

func startReceiver() {
	receiver.Run()
}

func newReceiverCmd() *cobra.Command {
	receiverCmd := &cobra.Command{
		Use:   "receiver",
		Short: "Execute tapirjr receiver",
		Run: func(cmd *cobra.Command, args []string) {
			startReceiver()
		},
	}

	receiverCmd.PersistentFlags().StringVarP(&receiver.ReceiverOpts.GobgpAddr, "gobgp-addr", "g", "127.0.0.1:50051", "specify the address and port for gobgp")
	receiverCmd.PersistentFlags().StringVarP(&receiver.ReceiverOpts.Port, "port", "p", "0.0.0.0:50051", "specify the port that taiprjr listens on")

	return receiverCmd
}
