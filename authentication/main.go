package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"assetra/authentication/auth"
	"assetra/authentication/repository"
	"assetra/db"
	"assetra/pb"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var (
	local bool
	port  int
)

func init() {
	flag.BoolVar(&local, "local", true, "run authentication service local")
	flag.IntVar(&port, "port", 9001, "port for the authentication service")
	flag.Parse()
}

func main() {

	if local {
		// TODO: Remove this and make K8s especial deployment
		err := godotenv.Load("../.env") // load .env file
		if err != nil {
			log.Panic("Failed to load .env file")
		}
	}

	// Initialize database connection
	cfg := db.NewConfig()
	conn := db.NewConnection(cfg)
	defer conn.Close()

	userRepository := repository.NewUserRepository(conn)
	refreshTokenRepository := repository.NewRefreshTokenRepository(conn)
	authService := auth.NewAuthService(userRepository, refreshTokenRepository)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Authentication service listening on port %d", port)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
