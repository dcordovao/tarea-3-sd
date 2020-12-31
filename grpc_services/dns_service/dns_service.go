package dns_service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

const zf_folder_path_1 = "servidor_dns/zf_files1"
const zf_folder_path_2 = "servidor_dns/zf_files2"
const zf_folder_path_3 = "servidor_dns/zf_files3"

var zf_folder_paths = []string{zf_folder_path_1, zf_folder_path_2, zf_folder_path_3}

type ClockVector struct {
	X, Y, Z int
}

type Server struct {
	Relojes map[string]ClockVector
}

// Esta funcion suma 1 en la posicion del vector correspondiente al server indicado por index
func sumar_uno_a_reloj(c ClockVector, index int) ClockVector {
	if index == 0 {
		return ClockVector{X: c.X + 1, Y: c.Y, Z: c.Z}
	} else if index == 1 {
		return ClockVector{X: c.X, Y: c.Y + 1, Z: c.Z}
	} else if index == 2 {
		return ClockVector{X: c.X, Y: c.Y, Z: c.Z + 1}
	} else {
		log.Fatal("ERROR, indice de reloj erroneo")
		return ClockVector{}
	}
}

// Funcion auxiliar para printear y debbuguear
func reloj_a_string(c ClockVector) string {
	return strconv.Itoa(c.X) + " " + strconv.Itoa(c.Y) + " " + strconv.Itoa(c.Z)
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

func (s *Server) CreateName(ctx context.Context, nombre *NewName) (*CommandResponse, error) {

	file_name := zf_folder_paths[nombre.IdDns] + "/" + nombre.Domain + ".zf"
	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		log.Printf("Creando dominio: " + nombre.Domain + ".zf")
		// Se crea en 0 el reloj vector de este dominio
		s.Relojes[nombre.Domain] = ClockVector{0, 0, 0}
		f, err := os.Create(file_name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		f.WriteString(nombre.Name + "." + nombre.Domain + " IN A " + nombre.Ip)
		// Actualizar reloj
		s.Relojes[nombre.Domain] = sumar_uno_a_reloj(s.Relojes[nombre.Domain], int(nombre.IdDns))
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
			// Actualizar reloj
			s.Relojes[nombre.Domain] = sumar_uno_a_reloj(s.Relojes[nombre.Domain], int(nombre.IdDns))
		} else {
			return &CommandResponse{Body: "Nombre no creado. El nombre ya estaba registrado en el dominio...", Clock: nil}, nil
		}
	}
	ultimo_reloj := s.Relojes[nombre.Domain]
	fmt.Println("Se creo con exito! Reloj dominio " + nombre.Domain + ": " + reloj_a_string(ultimo_reloj))
	reloj_mensaje := &ClockMessage{X: int64(ultimo_reloj.X), Y: int64(ultimo_reloj.Y), Z: int64(ultimo_reloj.Z)}
	response := &CommandResponse{Body: "Nombre creado con exito", Clock: reloj_mensaje}
	return response, nil

}

// Suponemos que al actualziar nombre, se da solo "nombre", y el dominio siempre se mantiene
func (s *Server) Update(ctx context.Context, update_info *UpdateInfo) (*CommandResponse, error) {
	file_name := zf_folder_paths[update_info.IdDns] + "/" + update_info.Domain + ".zf"
	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		log.Println("Se trato de hacer update a un dominio no existente...")
		return &CommandResponse{Body: "ERROR! Ese dominio no existe...", Clock: nil}, nil
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

		// Actualizar reloj
		s.Relojes[update_info.Domain] = sumar_uno_a_reloj(s.Relojes[update_info.Domain], int(update_info.IdDns))

		// Enviar reloj como respuesta
		ultimo_reloj := s.Relojes[update_info.Domain]
		fmt.Println("Se creo con exito! Reloj dominio " + update_info.Domain + ": " + reloj_a_string(ultimo_reloj))
		reloj_mensaje := &ClockMessage{X: int64(ultimo_reloj.X), Y: int64(ultimo_reloj.Y), Z: int64(ultimo_reloj.Z)}
		return &CommandResponse{Body: "Información actualizada con exito!", Clock: reloj_mensaje}, nil
	} else {
		return &CommandResponse{Body: "ERROR! No se encontro ese nombre en el dominio...", Clock: nil}, nil
	}

}

// Suponemos que al actualziar nombre, se da solo "nombre", y el dominio siempre se mantiene
func (s *Server) Delete(ctx context.Context, delete_info *DeleteInfo) (*CommandResponse, error) {
	file_name := zf_folder_paths[delete_info.IdDns] + "/" + delete_info.Domain + ".zf"
	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		log.Println("Se trato de hacer delete a un dominio no existente...")
		fmt.Printf(file_name + "\n")
		return &CommandResponse{Body: "ERROR! Ese dominio no existe...", Clock: nil}, nil
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
		if lname == delete_info.Name {
			index = i
			previusly_created = true
			continue
		}
	}

	if previusly_created {
		// Si se encontro la linea la eliminamos
		txtlines = append(txtlines[:index], txtlines[index+1:]...)

		// Eliminamos el archivo y lo creamos de nuevo sin la linea eliminada
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

		// Actualizar reloj
		s.Relojes[delete_info.Domain] = sumar_uno_a_reloj(s.Relojes[delete_info.Domain], int(delete_info.IdDns))

		// Enviar reloj como respuesta
		ultimo_reloj := s.Relojes[delete_info.Domain]
		fmt.Println("Se creo con exito! Reloj dominio " + delete_info.Domain + ": " + reloj_a_string(ultimo_reloj))
		reloj_mensaje := &ClockMessage{X: int64(ultimo_reloj.X), Y: int64(ultimo_reloj.Y), Z: int64(ultimo_reloj.Z)}
		return &CommandResponse{Body: "Información eliminada con exito!", Clock: reloj_mensaje}, nil
	} else {
		return &CommandResponse{Body: "ERROR! No se encontro ese nombre en el dominio...", Clock: nil}, nil
	}
}

// Esta funcion solo retorna el vector reloj del nombre.dominio solicitado
func (s *Server) GetClock(ctx context.Context, domain *Message) (*ClockMessage, error) {
	if val, ok := s.Relojes[domain.Body]; ok {
		return &ClockMessage{X: int64(val.X), Y: int64(val.Y), Z: int64(val.Z)}, nil
	} else {
		// Si no existe retornar reloj con 0,0,0
		return &ClockMessage{X: 0, Y: 0, Z: 0}, nil
	}
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
