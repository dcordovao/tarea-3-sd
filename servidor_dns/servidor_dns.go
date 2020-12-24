package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen on port 9001: %v", err)
	}

	fmt.Println("Escuchando en el puerto 9001...")

	// Setear server
	s := dns_service.Server{}

	grpcServer := grpc.NewServer()
	dns_service.RegisterDnsServiceServer(grpcServer, &s)

	////// Servicio de clientes ///////
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 9000: %v", err)
	}
	fmt.Println("Server corriendo...")
}
