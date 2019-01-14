package sender

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	ptapi "github.com/nstgt/tapirjr/api/proto"
	gobgpapi "github.com/osrg/gobgp/api"
	"google.golang.org/grpc"
)

var SenderOpts struct {
	GobgpAddr string
	PeerAddrs string
}

const (
	MAX_PATH_CAPACITY = 10000
)

var ac gobgpapi.GobgpApiClient

func Run() {
	opt := grpc.WithInsecure()
	conn, err := grpc.Dial(SenderOpts.GobgpAddr, opt)
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}

	pathChan := make(chan *gobgpapi.Path, MAX_PATH_CAPACITY)

	var senders []sender
	addrs := parsePeerAddrs(SenderOpts.PeerAddrs)
	for _, addr := range addrs {
		cc := getClientConn(addr)
		s := sender{address: addr, cliconn: cc}
		senders = append(senders, s)
	}

	ac = gobgpapi.NewGobgpApiClient(conn)
	// Only support for Afi: IPv4,IPv6 and Safi: Unicast
	go monitorRib(gobgpapi.Family_AFI_IP, gobgpapi.Family_SAFI_UNICAST, pathChan)
	go monitorRib(gobgpapi.Family_AFI_IP6, gobgpapi.Family_SAFI_UNICAST, pathChan)
	go distributePath(pathChan, senders)

}

func monitorRib(afi gobgpapi.Family_Afi, safi gobgpapi.Family_Safi, pathChan chan *gobgpapi.Path) {
	stream, err := ac.MonitorTable(context.Background(), &gobgpapi.MonitorTableRequest{
		TableType:  gobgpapi.TableType_ADJ_IN,
		Name:       "",
		Family:     &gobgpapi.Family{Afi: afi, Safi: safi},
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

func distributePath(pathChan chan *gobgpapi.Path, senders []sender) {
	for {
		path, ok := <-pathChan
		if !ok {
			break
		}

		if isPathIncludeNilSourceIPNeighborIP(path) {
			continue
		}

		for _, s := range senders {
			go s.sendPath(path)
		}
	}
}

func isPathIncludeNilSourceIPNeighborIP(path *gobgpapi.Path) (res bool) {
	p := *path
	if p.SourceId == "<nil>" {
		return true
	}
	if p.NeighborIp == "<nil>" {
		return true
	}
	return false
}

func parsePeerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}

func getClientConn(addr string) (cc *grpc.ClientConn) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return cc
}

type sender struct {
	address string
	cliconn *grpc.ClientConn
}

func (s sender) sendPath(path *gobgpapi.Path) {
	cli := ptapi.NewPathTransferClient(s.cliconn)

	// For debug
	fmt.Println(path)
	_, err := cli.Transmit(context.Background(), path)
	if err != nil {
		log.Fatal(err)
	}
}
