package dns_service

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
)

type Server struct {
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received a new message body from client: %s", message.Body)
	return &Message{Body: "Conectado desde name_service! "}, nil
}

func (s *Server) CreateName(ctx context.Context, nombre *NewName) (*Message, error) {

	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat("servidor_dns/zf_files/" + nombre.Domain); os.IsNotExist(err) {
		log.Printf("Creando dominio: " + nombre.Domain)

		f, err := os.Create("servidor_dns/zf_files/" + nombre.Domain)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Al crear el dominio seteamos su reloj en 0 0 0
		f.WriteString("0 0 0\n")
		f.WriteString(nombre.Name + "."+ nombre.Domain + " IN A " + nombre.Ip)
		defer f.Close()
	} else {
		// Leer el archivo, leer linea por linea, y si el nombre no existe es creado.
	}

	return &Message{Body: "Nombre creado con exito"}, nil

}
