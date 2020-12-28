package main

import (
	"log"

	"github.com/dcordova/sd_tarea3/grpc_services/broker_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"bufio"
	//"io"
	"os"
	"fmt"	
	"strings"
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

	//------------------------------------------------------------------------//
	//------------------ PEDIR COMANDO AL CLIENTE GET ------------------------//
	//------------------------------------------------------------------------//

	fmt.Println("")
	fmt.Println("---------------------")
	fmt.Println("Client Shell")
	fmt.Println("---------------------")

	for true {
		//// GET NOMBRE.DOMINIO ////

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("-> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		params := strings.Split(input, " ")
		
		if len(params) != 2 {
			fmt.Println("Cuidado!, comando get debería tener 1 parametro...\n")
			continue
		}

		message := broker_service.Message{
			Body: params[1],
		}
		response, err := s.EnviarDom(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}

		log.Printf("Response from Broker: %s", response.Body)
	}

}
