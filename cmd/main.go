package main

import (
	"auth/api"
	"auth/api/handler"
	"auth/config"
	pb "auth/genproto/user"
	"auth/logs"
	"auth/service"
	"auth/storage/postgres"
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Db *sql.DB

func main() {
	listener, err := net.Listen("tcp", config.Load().Server.USER_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	Db, err = postgres.ConnectionDb()
	if err != nil {
		log.Fatal(err)
	}

	logger := logs.NewLogger()
	service1 := service.NewUserService(Db, logger)
	service2 := service.NewNotificationsService(Db, logger)
	defer service1.User.Close()

	server := grpc.NewServer()
	pb.RegisterUserServer(server, service1)
	pb.RegisterNotificationsServer(server, service2)
	log.Printf("Server listening at %v", listener.Addr())
	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Server listening at %v", listener.Addr())

	hand := NewHandler()
	router := api.Router(hand)
	err = router.Run(config.Load().Server.USER_ROUTER)
	if err != nil {
		log.Fatal(err)
	}
}

func NewHandler() *handler.Handler {

	conn, err := grpc.NewClient(config.Load().Server.USER_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error while connecting authentication service ", err)
	}

	return &handler.Handler{
		User: pb.NewUserClient(conn),
		Log:  logs.NewLogger(),
	}
}
