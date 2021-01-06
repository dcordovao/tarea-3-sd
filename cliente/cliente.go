package main

import (
	"log"

	"bufio"

	"github.com/dcordova/sd_tarea3/grpc_services/broker_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	//"io"
	"fmt"
	"os"
	"strings"
	"strconv" // Conversion de strings a int y viceversa
)

type ClockVector struct {
	X, Y, Z int
}

// Reloj ya visto, junto con su ip
type clientSeenClock struct {
	vector ClockVector
	ip     string
	idDns  int
}

// funcion auxiliar
func clock_to_struct(a *broker_service.ClockMessage) ClockVector {
	return ClockVector{X: int(a.X), Y: int(a.Y), Z: int(a.Z)}
}

// Esta funcion entrega true, si el primer vector es igual o mas reciente que el segundo
func compare_clocks(new_clock ClockVector, last_clock ClockVector) bool {
	return new_clock.X >= last_clock.X && new_clock.Y >= last_clock.Y && new_clock.Z >= last_clock.Z
}

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

	var clientClocks map[string]clientSeenClock
	clientClocks = make(map[string]clientSeenClock)

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
			Body: params[1]+" ",
		}

		name := params[1]
		name_split := strings.Split(name, ".")

		// Primera vez, el resultado viene de una ip al azar 
		r, err := s.Connect(context.Background(), &message)
		random_id := strconv.Itoa(int(r.Iddns)+1)
		fmt.Println("El broker eligio el dns de id: "+random_id+" con la ip: "+r.Ipdns)

		
		// Si el broker no le achunta a la maquina hay que revisar el reloj
		// y pasar al broker la ip donde si esta el nombre.dominio
		if err != nil {			
			log.Fatalf("Error al conectarse al Connect del Broker: %s", err)
		} else {

			estado := strings.Split(r.Body , " ")
			tipo_error := estado[len(estado)-1]


			var ip_connection string
			var id_dns int
			var domain_clock ClockVector

			// Cuando escogio la ip donde estaba el nombre.dom
			if tipo_error != "Nombre" && tipo_error != "Dominio" {
				domain_clock = clock_to_struct(r.Clock)					
				ip_connection = strings.Split(r.Body, " ")[0]				
				id_dns = int(r.Iddns)
				fmt.Println(r.Body, "Reloj: [", r.Clock.X, r.Clock.Y, r.Clock.Z, "]")

				// MONOTONIC READ
				// Si es que existe un vector visto por el cliente en ese dominio	
				if latest_clock, ok := clientClocks[name_split[1]]; ok {
					if !compare_clocks(domain_clock, latest_clock.vector) {
						// Cambiar la conexion a la ultima ip vista
						ip_connection = latest_clock.ip
						id_dns = latest_clock.idDns

						message := broker_service.Message{
							Body: params[1] +" "+ ip_connection +" "+ strconv.Itoa(id_dns),
						}				

						// Segunda vez, el resultado viene de ip obtenida ultima vez en el vectr 
						fmt.Println("Esta version es anterior a la ultima leída!")						
						fmt.Println("Se cambio la conexion a la ultima ip vista para este dominio: " + strconv.Itoa(id_dns+1))//ip_connection)
						r, err = s.Connect(context.Background(), &message)
						fmt.Println(r.Body, "Reloj: [", r.Clock.X, r.Clock.Y, r.Clock.Z, "]")
					}	
				}	
			// Cuando falla debe revisar si está en otro dns		
			} else {
				fmt.Println(r.Body)
							
				if latest_clock, ok := clientClocks[name_split[1]]; ok {					
					
					ip_connection = latest_clock.ip
					id_dns = latest_clock.idDns

					message := broker_service.Message{
						Body: params[1] +" "+ ip_connection +" "+ strconv.Itoa(id_dns),
					}				

					// Segunda vez, el resultado viene de ip obtenida ultima vez en el vectr 					
					fmt.Println("Se cambio la conexion a la ultima ip vista para este dominio: " + strconv.Itoa(id_dns+1))
					r, err = s.Connect(context.Background(), &message)
					if r.Clock != nil {
						fmt.Println(r.Body, "[", r.Clock.X, r.Clock.Y, r.Clock.Z, "]")																											
					} else {
						fmt.Printf("El sitio "+params[1]+" ha sido eliminado")
					}
					
				}	
			}	
			if r.Clock != nil {		
				clientClocks[name_split[1]] = clientSeenClock{vector: clock_to_struct(r.Clock), ip: ip_connection, idDns: id_dns}				
			}
		}				
	}
}