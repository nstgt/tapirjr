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

var ac gobgpapi.GobgpApiClient

func Run() {
	conn, err := grpc.Dial(ServerOpts.GobgpAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ac = gobgpapi.NewGobgpApiClient(conn)

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

type server struct{}

func (s *server) Transmit(ctx context.Context, path *gobgpapi.Path) (*ptapi.PathTransferResponse, error) {
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
