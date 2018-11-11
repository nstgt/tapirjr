package server

import (
	"fmt"
	"io"
	"log"
	"net"

	ptapi "github.com/7310510/tapirjr/api/proto"
	"google.golang.org/grpc"
)

type server struct{}

func Run(gobgp_addr string, port string) {
	opt := grpc.WithInsecure()
	conn, err := grpc.Dial(gobgp_addr, opt)
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}

	lis, err := net.Listen("tcp", port)
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

func (s *server) Transmit(stream ptapi.PathTransfer_TransmitServer) error {
	for {
		path, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(path)
	}
	return nil
}
