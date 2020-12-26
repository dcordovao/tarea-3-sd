package dns_service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
)

type Server struct {
}

func check_if_name_in_domain(file_name string, new_name string) bool {
	previusly_created := false

	file, err := os.Open(file_name)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, eachline := range txtlines[1:] {
		//fmt.Println(eachline)
		lname := strings.Split(strings.Split(eachline, " ")[0], ".")[0]
		if lname == new_name {
			previusly_created = true
			continue
		}
	}
	return previusly_created
}

//////   Esta función era del tutorial pero la dejamos    ///////
//////   para ratificar la conexion con el servidor       ///////
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received a new message body from client: %s", message.Body)
	return &Message{Body: "Conectado desde name_service! "}, nil
}

func (s *Server) CreateName(ctx context.Context, nombre *NewName) (*Message, error) {

	file_name := "servidor_dns/zf_files/" + nombre.Domain
	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		log.Printf("Creando dominio: " + nombre.Domain)

		f, err := os.Create(file_name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Al crear el dominio seteamos su reloj en 0 0 0
		f.WriteString("0 0 0\n")
		f.WriteString(nombre.Name + "." + nombre.Domain + " IN A " + nombre.Ip)
		defer f.Close()
	} else {
		// Leer el archivo, leer linea por linea, y si el nombre no existe es creado.}

		previusly_created := check_if_name_in_domain(file_name, nombre.Name)

		// Si el nombre no existia en el dominio no existia, se puede crear crea
		if !previusly_created {
			f, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			if _, err = f.WriteString("\n" + nombre.Name + "." + nombre.Domain + " IN A " + nombre.Ip); err != nil {
				panic(err)
			}
		} else {
			return &Message{Body: "Nombre no creado. El nombre ya estaba registrado en el dominio..."}, nil
		}
	}

	return &Message{Body: "Nombre creado con exito"}, nil

}
