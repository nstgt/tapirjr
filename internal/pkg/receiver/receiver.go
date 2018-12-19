package receiver

import (
	"context"
	"fmt"
	"log"
	"net"

	ptapi "github.com/nstgt/tapirjr/api/proto"
	gobgpapi "github.com/osrg/gobgp/api"
	"google.golang.org/grpc"
)

var ReceiverOpts struct {
	GobgpAddr string
	Port      string
}

var ac gobgpapi.GobgpApiClient

func Run() {
	conn, err := grpc.Dial(ReceiverOpts.GobgpAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ac = gobgpapi.NewGobgpApiClient(conn)

	lis, err := net.Listen("tcp", ReceiverOpts.Port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	ptapi.RegisterPathTransferServer(s, new(receiver))
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

type receiver struct{}

func (r *receiver) Transmit(ctx context.Context, path *gobgpapi.Path) (*ptapi.PathTransferResponse, error) {
	// ForDebug
	fmt.Println(path)

	_, err := ac.AddPath(context.Background(), &gobgpapi.AddPathRequest{
		Path: path,
	})
	if err != nil {
		log.Fatal(err)
	}

	res := ptapi.PathTransferResponse{
		Status: "ok",
	}
	return &res, nil
}
