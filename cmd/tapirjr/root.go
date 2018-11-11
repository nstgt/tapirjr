package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var globalOpts struct {
	GobgpAddr string
	PeerAddrs string
	Port      string
}

var rootCmd = &cobra.Command{
	Use: "tapirjr",
	Run: func(c *cobra.Command, args []string) {},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globalOpts.GobgpAddr, "gobgp-addr", "g", "127.0.0.1:50051", "specify the address and port for gobgp (default: 127.0.0.1:50051")
	rootCmd.PersistentFlags().StringVarP(&globalOpts.PeerAddrs, "peer-addrs", "s", "", "address and port for server\nex) -s 10.0.0.1:50051,10.0.0.2:50051")
	rootCmd.PersistentFlags().StringVarP(&globalOpts.Port, "port", "p", ":50051", "specify the port that taiprjr listens on (default: :50051)")
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
