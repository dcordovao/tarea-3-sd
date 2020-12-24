package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dcordova/sd_tarea3/grpc_services/broker_service"
	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
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
		Body: "Conectandose desde admin...",
	}

	response, err := s.SayHello(context.Background(), &message)
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
		//fmt.Printf("%v", params)

		option := params[0]

		/// Create: create <nombre.dominio> <ip>
		if option == "create" {
			if len(params) != 3 {
				log.Printf("ERROR!, comando create debería tener 3 parametros...")
				continue
			}

			// Aqui comunicarse con el BROKER y obtener una ip de un dns

			ip_dns := ":9001"
			name := params[1]
			name_split := strings.Split(name, ".")
			if len(name_split) != 2 {
				log.Printf("ERROR! nombre.dominio mal formateado...")
				continue
			}

			// Aqui comprobar si la IP esta biem formateada

			// Enviar al servidor dns el nombre que se quiere crear
			new_name := dns_service.NewName{Name: name_split[0], Domain: name_split[1], Ip: params[2]}

			var conn_dns *grpc.ClientConn
			conn_dns, err := grpc.Dial(ip_dns, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn_dns.Close()

			s_dns := dns_service.NewDnsServiceClient(conn_dns)

			response, err := s_dns.CreateName(context.Background(), &new_name)
			if err != nil {
				log.Fatalf("Error al tratar de crear nombre: %s", err)
			}
			log.Printf("Response from Server: %s", response.Body)
		}

		/// OPCION 2:
		if option == "update" {

		}
	}
}
