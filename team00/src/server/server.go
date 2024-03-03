package main

import (
	"log"
	"math/rand"
	"net"
	"time"

	t "team00/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type customTransmitterServer struct {
	t.UnimplementedTransmitterServer
}

func (s *customTransmitterServer) ListRequests(_ *t.Request, stream t.Transmitter_ListRequestsServer) error {
	rd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	mean := -10.0 + rand.Float64()*20
	sd := 0.3 + rand.Float64()*(1.5-0.3)
	id := uuid.New().String()

	for {
		response := &t.Response{
			Uuid:      id,
			Frequency: rd.NormFloat64()*sd + mean,
			Timestamp: time.Now().Unix(),
		}

		err := stream.Send(response)
		log.Println(response, mean, sd)
		if err != nil {
			return err
		}

		time.Sleep(250 * time.Millisecond)
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:3333")

	if err != nil {
		log.Fatalf("Ошибка прослушивания порта: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	t.RegisterTransmitterServer(grpcServer, &customTransmitterServer{})

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}

}
