package main

import (
	"assetra/pb"
	"assetra/statesvc/resthandlers"
	"assetra/statesvc/routes"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	authAddr string
	local    bool
	port     int
)

func init() {
	flag.BoolVar(&local, "local", true, "run authentication service local")
	flag.IntVar(&port, "port", 9000, "port to run authentication API")
	flag.StringVar(&authAddr, "auth_addr", "localhost:9001", "authentication service adress")
	flag.Parse()
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// should read the credentials from the env vars and set the client credentials
	return nil, nil
}

func main() {

	var conn *grpc.ClientConn
	if local {
		// TODO: Remove this and make K8s especial deployment
		err := godotenv.Load("../.env") // load .env file
		if err != nil {
			log.Panic("Failed to load .env file")
		}
		conn, err = grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal("error creating the local dev gRPC Client", err)
		}
	} else {

		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("error loading tls credentials", err)
		}

		conn, err = grpc.NewClient(authAddr, grpc.WithTransportCredentials(tlsCredentials))
		if err != nil {
			log.Fatal("error creating the gRPC Client", err)
		}
	}

	client := pb.NewAuthServiceClient(conn)
	authHandlers := resthandlers.NewAuthHandlers(client)
	authRoutes := routes.NewAuthRoutes(authHandlers)

	router := mux.NewRouter().StrictSlash(true)
	routes.Install(router, authRoutes)

	log.Printf("API service running on [::]:%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routes.WithCORS(router)))
}
