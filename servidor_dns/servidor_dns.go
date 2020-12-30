package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"google.golang.org/grpc"
)

func main() {

	// Se le puede dar un parametro al server que indica su indice
	// Esto no se usara en la version final, ya que todos los servers seran equivalentes
	fmt.Println("Len(Args): " + strconv.Itoa(len(os.Args)))

	var puerto string
	if len(os.Args) == 2 {
		puerto = ":900" + os.Args[1]
	} else {
		puerto = ":9001"
	}

	lis, err := net.Listen("tcp", puerto)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", puerto, err)
	}

	fmt.Println("Escuchando en el puerto " + puerto + "...")

	// Setear server
	s := dns_service.Server{}

	grpcServer := grpc.NewServer()
	dns_service.RegisterDnsServiceServer(grpcServer, &s)

	////// Servicio de clientes ///////
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server on port %s: %v", puerto, err)
	}
	fmt.Println("Server corriendo...")
}
