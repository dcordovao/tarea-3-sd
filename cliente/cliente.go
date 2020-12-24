package main

import (
	"log"

	"github.com/dcordova/sd_tarea3/grpc_services/broker_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"bufio"
	//"io"
	//"os"
	//"strconv" // Conversion de strings a int y viceversa
)

func main() {

	//------------------------------------------------------
	//////// Conectarse como cliente al NameService ////////
	//------------------------------------------------------
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	s := broker_service.NewBrokerServiceClient(conn)

	// Hello world
	message := broker_service.Message{
		Body: "Conectandose desde cliente...",
	}

	response, err := s.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)
}
