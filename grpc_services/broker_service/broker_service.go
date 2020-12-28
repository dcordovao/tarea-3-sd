package broker_service

import (
	"log"	
	"github.com/dcordova/sd_tarea3/grpc_services/dns_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received a new message body from client: %s", message.Body)
	return &Message{Body: "Saludos desde broker_server! "}, nil
}

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