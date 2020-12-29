package broker_service

import (
	"log"	
	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"	
	//"strings"
	"math/rand"
)


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
func (s *Server) EnviarDom(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received from client: %s, now sending to dns sevice", message.Body)


	//----------------------------------------------------------//
	//----------- EN ESTA PARTE SE PIDE AL DNS LA --------------//
	//----------- IP Y EL RELOJ SEGUN EL DOMINIO SOLICITADO ----//
	//----------------------------------------------------------//	

	var conn_dns *grpc.ClientConn

	//IP_DNS = ":9001" Esto cambia segun quien lo tiene
	conn_dns, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn_dns.Close()

	// Mensage que llego desde el cliente
	get_request := dns_service.Message{Body: message.Body}

	s_dns := dns_service.NewDnsServiceClient(conn_dns)

	response, err := s_dns.GetName(context.Background(), &get_request)
	if err != nil {
		log.Fatalf("Error al tratar de crear nombre: %s", err)
	}
	log.Printf("Response from Server: %s", response.Body)

	return &Message{Body: response.Body}, nil
}

//////   Recibe Verbo     ///////
//////   Retorna IP       ///////
func (s *Server) EnviarVerbo(ctx context.Context, operacion *Message) (*Message, error) {
	
	IPs := []string{"10.10.28.121", "10.10.28.122", "10.10.28.123"}
     // Seleccion random de servidor
    //randomDNS := IPs[rand.Intn(len(IPs))]    
    randomDNS := []string{"1/", "2/", "3/"}[rand.Intn(len(IPs))]			
	
	return &Message{Body: randomDNS}, nil//\n Respuesta: " +response.Body}, nil
}