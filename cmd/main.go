package main

import (
	"context"
	"log"
	"net"

	"github.com/justIGreK/Reminders/cmd/handler"
	"github.com/justIGreK/Reminders/internal/repository"
	"github.com/justIGreK/Reminders/internal/service"
	"github.com/justIGreK/Reminders/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	tz, err:= client.NewTimezoneClient("localhost:50055")
	if err != nil {
		log.Fatal(err)
	}
	db := repository.CreateMongoClient(ctx)
	rmsRepo := repository.NewRemindersRepository(db)
	rmsSRV := service.NewRemindersService(rmsRepo, tz)
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	handler := handler.NewHandler(grpcServer, rmsSRV)
	handler.RegisterServices()
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
