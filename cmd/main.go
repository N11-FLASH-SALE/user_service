package main

import (
	"auth/api"
	"auth/api/handler"
	"auth/genproto/user"
	"auth/pkg/logger"
	"auth/service"
	"auth/storage/postgres"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db, err := postgres.ConnectionDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", "auth:50051") 
	if err != nil {
		log.Fatalf("error while listening: %v", err)
	}
	defer lis.Close()

	serviceUser := service.NewUserService(db, logger.NewLogger())
	server := grpc.NewServer()
	user.RegisterUsersServer(server, serviceUser)
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		err = server.Serve(lis)
		if err != nil {
			log.Fatalf("error while serving: %v", err)
		}
	}()
	hand := NewHandler()
	router := api.Router(hand)
	log.Println("server is running")
	log.Fatal(router.Run("auth:8085"))
}

func NewHandler() *handler.Handler {
	conn, err := grpc.NewClient("auth:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	return &handler.Handler{User: user.NewUsersClient(conn), Log: logger.NewLogger()}
}
