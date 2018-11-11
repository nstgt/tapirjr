package main

import (
	"github.com/7310510/tapirjr/internal/pkg/client"
	"github.com/7310510/tapirjr/internal/pkg/server"
)

func main() {
	execute()

	go server.Run(globalOpts.GobgpAddr, globalOpts.Port)
	go client.Run(globalOpts.GobgpAddr, globalOpts.PeerAddrs)

	for {
	}
}
