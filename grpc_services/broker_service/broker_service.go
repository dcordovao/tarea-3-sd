package broker_service

import (
	"log"
	//"fmt"
	"strconv"

	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"math/rand"
	"strings"
)

const ip_dns_1 = ":9001" //"10.10.28.121:9000"
const ip_dns_2 = ":9002" //"10.10.28.122:9000"
const ip_dns_3 = ":9003" //"10.10.28.123:9000"
var ips_dns = [...]string{ip_dns_1, ip_dns_2, ip_dns_3}

type Server struct {
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received a new message body from client: %s", message.Body)
	return &Message{Body: "Saludos desde broker_server! "}, nil
}

//////   Esta función recibe la request del cliente    ///////
//////   y retorna la IP y el Reloj                    ///////
func (s *Server) Connect(ctx context.Context, message *Message) (*CommandResponse, error) {
	log.Printf("Received from client: %s, now sending to dns sevice", message.Body)

	//----------------------------------------------------------//
	//----------- EN ESTA PARTE SE PIDE AL DNS LA --------------//
	//----------- IP Y EL RELOJ SEGUN EL DOMINIO SOLICITADO ----//
	//----------------------------------------------------------//
	mens := strings.Split(message.Body, " ")[0]

	name := strings.Split(mens, ".")[0]
	domain := strings.Split(mens, ".")[1]

	var random_ip string
	var random_int int
	if len(strings.Split(message.Body, " ")) > 2 {
		// Ya no es random
		//random_ip = strings.Split(message.Body, " ")[1]
		random_int, _ = strconv.Atoi(strings.Split(message.Body, " ")[2])

	} else {
		random_int = rand.Intn(len(ips_dns))
	}
	random_ip = ips_dns[random_int]

	var conn_dns *grpc.ClientConn

	//IP_DNS = ":9001" Esto cambia segun quien lo tiene
	conn_dns, err := grpc.Dial(random_ip, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn_dns.Close()

	// Mensage que llego desde el cliente
	get_request := dns_service.NewName{Name: name, Domain: domain, Ip: random_ip, IdDns: int64(random_int)}

	s_dns := dns_service.NewDnsServiceClient(conn_dns)

	response, err := s_dns.GetName(context.Background(), &get_request)
	if err != nil {
		log.Fatalf("Error al conectarse al GetName del Dns: %s", err)
	}
	estado := strings.Split(response.Body, " ")
	tipo_error := estado[len(estado)-1]

	if tipo_error != "Nombre" && tipo_error != "Dominio" {
		ip_id := response.Body //+ " " + strconv.Itoa(random_int)
		clock_res := &ClockMessage{X: response.Clock.X, Y: response.Clock.Y, Z: response.Clock.Z}
		return &CommandResponse{Body: ip_id, Clock: clock_res, Ipdns: random_ip, Iddns: int64(random_int)}, nil
	}
	// No estaba
	return &CommandResponse{Body: response.Body, Clock: nil, Ipdns: random_ip, Iddns: int64(random_int)}, nil
}

//////   Recibe Verbo     ///////
//////   Retorna IP       ///////
func (s *Server) EnviarVerbo(ctx context.Context, operacion *Message) (*DnsAddress, error) {

	// Seleccion random de servidor
	//randomDNS := IPs[rand.Intn(len(IPs))]
	random_int := rand.Intn(len(ips_dns))
	random_ip := ips_dns[random_int]

	return &DnsAddress{Ip: random_ip, IdDns: int64(random_int)}, nil
}

func (s *Server) PropagarCambios(ctx context.Context, message *Message) (*Message, error) {
	//var cambios_server_1 map[string][]string
	//cambios_server_1 = make(map[string][]string)

	// Siempre se conectara primero al dns 1
	//ip_dns = ip_dns_1
	//var modificaciones = dns_service.Modificaciones

	log.Println("Propagando cambios...")

	var conn_dns2 *grpc.ClientConn
	//IP_DNS = ":9001" Esto cambia segun quien lo tiene
	conn_dns2, err2 := grpc.Dial(ip_dns_2, grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("Could not connect: %s", err2)
	}
	defer conn_dns2.Close()

	// Mensage que llego desde el cliente

	s_dns2 := dns_service.NewDnsServiceClient(conn_dns2)

	id_dns2 := dns_service.IdDns{IdDns: 1, IpDns: ip_dns_1}
	response2, err2 := s_dns2.PropagarCambios(context.Background(), &id_dns2)
	if err2 != nil {
		log.Fatalf("Error al pedir al servidor 2 que propague camibios: %s", err2)
	}
	log.Println("Respuesta servidor 2:", response2.Body)

	var conn_dns3 *grpc.ClientConn
	//IP_DNS = ":9001" Esto cambia segun quien lo tiene
	conn_dns3, err3 := grpc.Dial(ip_dns_3, grpc.WithInsecure())
	if err3 != nil {
		log.Fatalf("Could not connect: %s", err3)
	}
	defer conn_dns3.Close()

	// Mensage que llego desde el cliente

	s_dns3 := dns_service.NewDnsServiceClient(conn_dns3)

	id_dns3 := dns_service.IdDns{IdDns: 2, IpDns: ip_dns_1}
	response3, err3 := s_dns3.PropagarCambios(context.Background(), &id_dns3)
	if err3 != nil {
		log.Fatalf("Error al pedir al servidor 3 que propague camibios: %s", err3)
	}
	log.Println("Respuesta servidor 3:", response3.Body)

	////////// PEDIR AL SERVIDOR 1 QUE MANDE LOS ZF A LOS SERVIDORES 2 y 3
	var conn_dns1 *grpc.ClientConn
	//IP_DNS = ":9001" Esto cambia segun quien lo tiene
	conn_dns1, err1 := grpc.Dial(ip_dns_1, grpc.WithInsecure())
	if err1 != nil {
		log.Fatalf("Could not connect: %s", err1)
	}
	defer conn_dns1.Close()

	// Mensage que llego desde el cliente

	s_dns1 := dns_service.NewDnsServiceClient(conn_dns1)

	target_ips := dns_service.TargetIps{IdDns: 0, Ip1: ip_dns_1, Ip2: ip_dns_2}
	response1, err1 := s_dns1.PropagarZfs(context.Background(), &target_ips)
	if err1 != nil {
		log.Fatalf("Error al pedir al servidor 1 envie los zf: %s", err3)
	}
	log.Println("Respuesta servidor 3:", response1.Body)

	// Pedir a los server que borren sus logs
	id1 := dns_service.IdDns{IdDns: 0, IpDns: ip_dns_1}
	r1, err1 := s_dns1.EliminarLogs(context.Background(), &id1)
	if err1 != nil {
		log.Fatalf("Error al pedir al servidor 1 envie los zf: %s", err1)
	}
	log.Println("Respuesta servidor 3:", r1.Body)

	id2 := dns_service.IdDns{IdDns: 1, IpDns: ip_dns_2}
	r2, err2 := s_dns2.EliminarLogs(context.Background(), &id2)
	if err2 != nil {
		log.Fatalf("Error al pedir al servidor 1 envie los zf: %s", err2)
	}
	log.Println("Respuesta servidor 3:", r2.Body)

	id3 := dns_service.IdDns{IdDns: 2, IpDns: ip_dns_3}
	r3, err2 := s_dns3.EliminarLogs(context.Background(), &id3)
	if err3 != nil {
		log.Fatalf("Error al pedir al servidor 1 envie los zf: %s", err3)
	}
	log.Println("Respuesta servidor 3:", r3.Body)

	return &Message{Body: "Saludos desde broker_server! "}, nil
}
