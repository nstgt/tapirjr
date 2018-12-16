package server

import (
	"context"
	"fmt"
	"log"
	"net"

	ptapi "github.com/7310510/tapirjr/api/proto"
	gobgpapi "github.com/osrg/gobgp/api"
	"google.golang.org/grpc"
)

var ServerOpts struct {
	GobgpAddr string
	Port      string
}

type server struct{}

func Run() {
	lis, err := net.Listen("tcp", ServerOpts.Port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	ptapi.RegisterPathTransferServer(s, new(server))
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *server) Transmit(ctx context.Context, path *gobgpapi.Path) (*ptapi.PathTransferResponse, error) {
	// WIP
	fmt.Println(path)
	res := ptapi.PathTransferResponse{
		Status: "ok",
	}
	return &res, nil
}
