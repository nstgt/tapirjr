package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	ptapi "github.com/7310510/tapirjr/api/proto"
	gobgpapi "github.com/osrg/gobgp/api"
	"google.golang.org/grpc"
)

const (
	MAX_PATH_CAPACITY = 10000
)

var ac gobgpapi.GobgpApiClient

func Run(gobgp_addr string, peer_addrs string) {
	opt := grpc.WithInsecure()
	conn, err := grpc.Dial(gobgp_addr, opt)
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}

	var clients []client

	addrs := parsePeerAddrs(peer_addrs)
	for _, addr := range addrs {
		c := client{address: addr}
		c.start()
		clients = append(clients, c)
	}

	pathChan := make(chan *gobgpapi.Path, MAX_PATH_CAPACITY)

	ac = gobgpapi.NewGobgpApiClient(conn)
	go monitorRib(pathChan)
	go distributePath(pathChan, clients)
}

func monitorRib(pathChan chan *gobgpapi.Path) {
	stream, err := ac.MonitorTable(context.Background(), &gobgpapi.MonitorTableRequest{
		Type:       gobgpapi.Resource_GLOBAL,
		Name:       "",
		Family:     &gobgpapi.Family{Afi: gobgpapi.Family_AFI_IP, Safi: gobgpapi.Family_SAFI_UNICAST},
		Current:    true,
		PostPolicy: true,
	})
	if err != nil {
		log.Fatalf("RPC error...: %v", err)
	}

	for {
		p, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("RPC error: %v", err)
		}

		path := p.Path
		pathChan <- path
	}
}

func distributePath(pathChan chan *gobgpapi.Path, clients []client) {
	for {
		path := <-pathChan
		for _, cli := range clients {
			go cli.sendPath(path)
		}
	}
}

func parsePeerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}

type client struct {
	address string
	stream  ptapi.PathTransfer_TransmitClient
}

func (c client) start() {
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("hoge")
	fmt.Println(conn.GetState())

	cli := ptapi.NewPathTransferClient(conn)
	stream, err := cli.Transmit(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	c.stream = stream
}

func (c client) sendPath(path *gobgpapi.Path) {
	fmt.Println(path)
	if err := c.stream.Send(path); err != nil {
		log.Fatal(err)
	}
}
