package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dcordova/sd_tarea3/grpc_services/broker_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	fmt.Println("Escuchando en el puerto 9000...")

	// Setear server
	s := broker_service.Server{}

	grpcServer := grpc.NewServer()
	broker_service.RegisterBrokerServiceServer(grpcServer, &s)

	////// Servicio de clientes ///////
	go func() {

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()

		broker := broker_service.NewBrokerServiceClient(conn)
		message := broker_service.Message{Body: "propagate"}
		for true {
			// Dormir hasta realizar propagacion
			time.Sleep(40 * time.Second)
			response, err := broker.PropagarCambios(context.Background(), &message)
			if err != nil {
				log.Fatalf("Error propagar cambios: %s", err)
			}
			log.Printf("Response from Server: %s", response.Body)
		}
	}()
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 9000: %v", err)
	}
	fmt.Println("Server corriendo...")
}
