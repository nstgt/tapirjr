package receiver

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	l, err := net.Listen("tcp", ReceiverOpts.Port)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	var wg sync.WaitGroup

	ctx, shutdown := context.WithCancel(context.Background())

	log.Println("receiver start...")
	wg.Add(1)
	go startServer(ctx, l, &wg)

	s := <-sigChan
	switch s {
	case syscall.SIGINT:
		log.Println("receiver shutdown...")
		shutdown()
		wg.Wait()
		log.Println("receiver shutdown completed, bye!")
	}
}

func startServer(ctx context.Context, l net.Listener, wg *sync.WaitGroup) {
	defer func() {
		l.Close()
		wg.Done()
	}()

	s := grpc.NewServer()
	ptapi.RegisterPathTransferServer(s, new(receiver))

	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		s.Stop()
		return
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
