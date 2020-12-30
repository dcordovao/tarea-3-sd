package dns_service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
)

const zf_folder_path_1 = "servidor_dns/zf_files1"
const zf_folder_path_2 = "servidor_dns/zf_files2"
const zf_folder_path_3 = "servidor_dns/zf_files3"

var zf_folder_paths = []string{zf_folder_path_1, zf_folder_path_2, zf_folder_path_3}

type ClockVector struct {
	x, y, z int
}

type Server struct {
	relojes map[string]ClockVector
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

	for _, eachline := range txtlines[:] {
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

	file_name := zf_folder_paths[nombre.IdDns] + "/" + nombre.Domain + ".zf"
	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		log.Printf("Creando dominio: " + nombre.Domain + ".zf")

		f, err := os.Create(file_name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()
		f.WriteString(nombre.Name + "." + nombre.Domain + " IN A " + nombre.Ip)

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

// Suponemos que al actualziar nombre, se da solo "nombre", y el dominio siempre se mantiene
func (s *Server) Update(ctx context.Context, update_info *UpdateInfo) (*Message, error) {
	file_name := zf_folder_paths[update_info.IdDns] + "/" + update_info.Domain + ".zf"
	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		return &Message{Body: "ERROR! Ese dominio no existe..."}, nil
	}
	// Si ya existe el archivo, lo leemos en busca del nombre que queremos modificar
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

	// Leer linea por linea para buscar un nombre que coincida
	previusly_created := false
	var index int
	for i, eachline := range txtlines[:] {
		//fmt.Println(eachline)
		lname := strings.Split(strings.Split(eachline, " ")[0], ".")[0]
		if lname == update_info.Name {
			index = i
			previusly_created = true
			continue
		}
	}

	// Si se encontro la linea la modificamos
	if previusly_created {
		if update_info.Opt == "name" {
			antigua_ip := strings.Split(txtlines[index], " ")[3]
			new_line := update_info.Value + "." + update_info.Domain + " IN A " + antigua_ip
			txtlines[index] = new_line
		} else {
			new_line := update_info.Name + "." + update_info.Domain + " IN A " + update_info.Value
			txtlines[index] = new_line
		}
		// Luego de extraer la linea y modificarla, borramos el arhivo y lo escribimos denuevo

		e := os.Remove(file_name)
		if e != nil {
			log.Fatal(e)
		}

		f, err := os.Create(file_name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		n_lines := len(txtlines)
		for _, eachline := range txtlines[:n_lines-1] {
			f.WriteString(eachline + "\n")
		}
		f.WriteString(txtlines[n_lines-1])
	} else {
		return &Message{Body: "ERROR! No se encontro ese nombre en el dominio..."}, nil
	}

	return &Message{Body: "Información actualizada con exito!"}, nil
}

func (s *Server) GetName(ctx context.Context, message *Message) (*Message, error) {
	//-----------------------------------------------------------//
	//----------- EN ESTA PARTE SE BUSCA EL DNS -----------------//
	//----------- CON EL DOMINIO SOLICITADO Y SE RETORNA LA IP --//
	//----------- Y EL RELOJ VECTORIAL ASOCIADO -----------------//
	//-----------------------------------------------------------//

	// Por ahora solo se retorna un string para probar. Borrar despues!!
	return &Message{Body: "Llego al DNS!!"}, nil
}
