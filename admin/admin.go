package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/dcordova/sd_tarea3/grpc_services/broker_service"
	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	//"bufio"
	//"io"
	//"os"
	"strconv" // Conversion de strings a int y viceversa
)

// Reloj ya visto, junto con su ip
type SeenClock struct {
	vector dns_service.ClockVector
	ip     string
	idDns  int
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

// funcion auxiliar
func clock_to_struct(a *dns_service.ClockMessage) dns_service.ClockVector {
	return dns_service.ClockVector{X: int(a.X), Y: int(a.Y), Z: int(a.Z)}
}

// Esta funcion entrega true, si el primer vector es igual o mas reciente que el segundo
func compare_clocks(new_clock dns_service.ClockVector, last_clock dns_service.ClockVector) bool {
	return new_clock.X >= last_clock.X && new_clock.Y >= last_clock.Y && new_clock.Z >= last_clock.Z
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

	var adminClocks map[string]SeenClock
	adminClocks = make(map[string]SeenClock)

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

		/// validar opcion 1
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
				log.Printf("Fallo en checkIPAddress: " + strconv.FormatBool(booleano))
				continue
			}
		}

		/// validar opcion 2
		if option == "update" {
			if len(params) != 4 {
				log.Printf("Cuidado!, comando update debería tener 3 parametros...")
				continue
			}
		}

		response, err := broker.EnviarVerbo(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error al enviar verbo: %s", err)
		}
		id_dns := int(response.IdDns)
		ip_connection := response.Ip

		log.Printf("El broker escogió el DNS %v la IP: %s", id_dns, response.Ip)

		//ip_dns := ":9001" //randomDNS
		var conn_dns *grpc.ClientConn
		conn_dns, err = grpc.Dial(ip_connection, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn_dns.Close()

		// Extraer nombre y dominio
		name := params[1]
		name_split := strings.Split(name, ".")

		// Se conecta a la primera IP, pero luego puede cambiar dependiendo del vector reloj que reciba
		s_dns := dns_service.NewDnsServiceClient(conn_dns)

		// READ YOUR WRITES
		// Obtener vector del dominio
		response_clock, err := s_dns.GetClock(context.Background(), &dns_service.Message{Body: name_split[1]})
		if err != nil {
			log.Fatalf("Error al tratar de obtener el reloj de dominio: %s", name_split[1])
		}
		domain_clock := clock_to_struct(response_clock)

		// READ YOUR WRITES
		// Si es que existe un vector visto por el admin de ese dominio
		if latest_clock, ok := adminClocks[name_split[1]]; ok {
			if !compare_clocks(domain_clock, latest_clock.vector) {
				// Cambiar la conexion a la ultima ip vista
				ip_connection = latest_clock.ip
				conn_dns, err = grpc.Dial(ip_connection, grpc.WithInsecure())
				if err != nil {
					log.Fatalf("Could not connect: %s", err)
				}
				defer conn_dns.Close()
				s_dns = dns_service.NewDnsServiceClient(conn_dns)
				id_dns = latest_clock.idDns
				fmt.Println("Se cambio la conexion a la ultima ip vista: " + ip_connection)
			}
		}

		/// Create: create <nombre.dominio> <ip>
		if option == "create" {

			// Aqui comunicarse con el BROKER y obtener una ip de un dns
			new_ip := params[2]

			// Enviar al servidor dns el nombre que se quiere crear
			new_name := dns_service.NewName{Name: name_split[0], Domain: name_split[1], Ip: new_ip, IdDns: int64(id_dns)}

			// Pedir el vector reloj del dominio

			response, err := s_dns.CreateName(context.Background(), &new_name)
			if err != nil {
				log.Fatalf("Error al tratar de crear nombre: %s", err)
			}
			log.Printf("Response from Server: %s", response.Body)
			// Si el reloj recibido no es nulo, se guarda como ultimo reloj de ese dominio
			if response.Clock != nil {
				adminClocks[new_name.Domain] = SeenClock{vector: clock_to_struct(response.Clock), ip: ip_connection, idDns: id_dns}
			}
		}

		/// OPCION 2:
		if option == "update" {

			update_info := dns_service.UpdateInfo{Name: name_split[0], Domain: name_split[1], Opt: params[2], Value: params[3], IdDns: int64(id_dns)}

			response, err := s_dns.Update(context.Background(), &update_info)
			if err != nil {
				log.Fatalf("Error al tratar de crear nombre: %s", err)
			}
			log.Printf("Response from Server: %s", response.Body)
		}
	}
}
