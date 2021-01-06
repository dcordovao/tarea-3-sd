package broker_service

import (
	"log"
	//"fmt"
	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"

	"strings"
	"math/rand"
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
		random_int,_ = strconv.Atoi(strings.Split(message.Body, " ")[2])		
		
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
	estado := strings.Split(response.Body , " ")
	tipo_error := estado[len(estado)-1]

	if tipo_error != "Nombre" && tipo_error != "Dominio" {
		ip_id := response.Body + " " + strconv.Itoa(random_int)
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
