package main

import (
	"github.com/7310510/tapirjr/internal/pkg/client"
	"github.com/spf13/cobra"
)

func startClient() {
	go client.Run()
	for {
	}
}

func newClientCmd() *cobra.Command {
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "execute tapirjr client",
		Run: func(cmd *cobra.Command, args []string) {
			startClient()
		},
	}

	clientCmd.PersistentFlags().StringVarP(&client.ClientOpts.GobgpAddr, "gobgp-addr", "g", "127.0.0.1:50051", "specify the address and port for gobgp")
	clientCmd.PersistentFlags().StringVarP(&client.ClientOpts.PeerAddrs, "peer-addrs", "s", "", "address and port for server\nex) -s 10.0.0.1:50051,10.0.0.2:50051")

	return clientCmd
}
