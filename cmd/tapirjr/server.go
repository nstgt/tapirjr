package main

import (
	"github.com/7310510/tapirjr/internal/pkg/server"
	"github.com/spf13/cobra"
)

func startServer() {
	go server.Run()
	for {
	}
}

func newServerCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "execute tapirjr server",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	serverCmd.PersistentFlags().StringVarP(&server.ServerOpts.GobgpAddr, "gobgp-addr", "g", "127.0.0.1:50051", "specify the address and port for gobgp")
	serverCmd.PersistentFlags().StringVarP(&server.ServerOpts.Port, "port", "p", "0.0.0.0:50051", "specify the port that taiprjr listens on")

	return serverCmd
}
