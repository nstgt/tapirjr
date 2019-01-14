package main

import (
	"github.com/nstgt/tapirjr/internal/pkg/sender"
	"github.com/spf13/cobra"
)

func startSender() {
	sender.Run()
}

func newSenderCmd() *cobra.Command {
	senderCmd := &cobra.Command{
		Use:   "sender",
		Short: "Execute tapirjr sender",
		Run: func(cmd *cobra.Command, args []string) {
			startSender()
		},
	}

	senderCmd.PersistentFlags().StringVarP(&sender.SenderOpts.GobgpAddr, "gobgp-addr", "g", "127.0.0.1:50051", "specify the address and port for gobgp")
	senderCmd.PersistentFlags().StringVarP(&sender.SenderOpts.PeerAddrs, "peer-addrs", "s", "", "address and port for server\nex) -s 10.0.0.1:50051,10.0.0.2:50051")

	return senderCmd
}
