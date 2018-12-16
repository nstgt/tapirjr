package sender

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
	go monitorRib(pathChan)
	// WIP: only one connection can be used to send path
	go senders[0].sendPath(pathChan)

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

// WIP
//func distributePath(pathChan chan *gobgpapi.Path, senders []sender) {
//	for {
//		path, ok := <-pathChan
//		if !ok {
//			break
//		}
//		p := *path
//		for _, cli := range senders {
//			go cli.sendPath(p)
//		}
//		c := senders[0]
//		c.sendPath(p)
//	}
//}

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

func (s *sender) sendPath(pathChan chan *gobgpapi.Path) {
	//conn, err := grpc.Dial(s.address, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer conn.Close()

	cli := ptapi.NewPathTransferClient(s.cliconn)

	for {
		path, ok := <-pathChan
		if !ok {
			break
		}
		// WIP
		fmt.Println(path)
		_, err := cli.Transmit(context.Background(), path)
		if err != nil {
			log.Fatal(err)
		}
	}
}
