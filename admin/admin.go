package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/dcordova/sd_tarea3/grpc_services/broker_service"
	//"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"bufio"
	//"io"
	//"os"
	"strconv" // Conversion de strings a int y viceversa
)

func verificar_nombre() {

}

// Validar si una IP es valida
func checkIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalida\n", ip)
		return false
	} else {
		//fmt.Printf("IP Address: %s - Valida\n", ip)
		return true
	}
}

func main() {

	//------------------------------------------------------
	//////// Conectarse como admin ////////
	//------------------------------------------------------
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	broker := broker_service.NewBrokerServiceClient(conn)

	// Hello world
	message := broker_service.Message{
		Body: "Conectandose desde admin...",
	}

	response, err := broker.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)

	fmt.Println("")
	fmt.Println("---------------------")
	fmt.Println("Admin Shell")
	fmt.Println("---------------------")

	for true {
		//// OPCIONES A ELEGIR ////

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("-> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)	

		params := strings.Split(input, " ")

		option := params[0]	

		message := broker_service.Message{
			Body: input,
		}

		/// Create: create <nombre.dominio> <ip>
		if option == "create" {
			if len(params) != 3 {							
				log.Printf("Cuidado!, comando create debería tener 2 parametros...")
				continue
			}

			// Aqui comunicarse con el BROKER y obtener una ip de un dns		
			name := params[1]
			name_split := strings.Split(name, ".")
			if len(name_split) != 2 {			
				log.Printf("Cuidado! nombre.dominio mal formateado...")
				continue
			}

			new_ip := params[2]

			booleano := checkIPAddress(new_ip)

			if !booleano {
				log.Printf("Fallo en checkIPAddress: "+strconv.FormatBool(booleano))
				continue				
			}							
		}

		/// OPCION 2:
		if option == "update" {
			if len(params) != 4 {			
				log.Printf("Cuidado!, comando update debería tener 3 parametros...")
				continue
			}										
		}

		response, err := broker.EnviarVerbo(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}

		log.Printf("Response from Server: %s", response.Body)
	}
}
