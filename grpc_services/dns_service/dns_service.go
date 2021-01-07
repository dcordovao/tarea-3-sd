package dns_service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const zf_folder_path_1 = "servidor_dns/zf_files"
const zf_folder_path_2 = "servidor_dns/zf_files"
const zf_folder_path_3 = "servidor_dns/zf_files"

var zf_folder_paths = []string{zf_folder_path_1, zf_folder_path_2, zf_folder_path_3}

type ClockVector struct {
	X, Y, Z int
}

type Server struct {
	Relojes map[string]ClockVector
}

// Esta funcion escribe en el log de un dominio, si no está, lo crea
func DomainLog(Name string, Domain string, Ip string, op string, IdDns int64) {
	file_name := zf_folder_paths[IdDns] + "/" + Domain + ".log"
	// Chequear si el log del dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		log.Printf("Creando log: " + Domain + ".log")

		f, err := os.Create(file_name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		f.WriteString(op + " " + Name + "." + Domain + " " + Ip)

		// Si ya existe, entonces agrega el comando al log
	} else {

		f, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString("\n" + op + " " + Name + "." + Domain + " " + Ip); err != nil {
			panic(err)
		}
	}
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
			var salto string
			if len(txtlines) == 0 {
				salto = ""
			} else {
				salto = "\n"
			}

			f, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}

			if _, err = f.WriteString(salto + nombre.Name + "." + nombre.Domain + " IN A " + nombre.Ip); err != nil {
				panic(err)
			}

			f.Close()
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

	DomainLog(nombre.Name, nombre.Domain, nombre.Ip, "create", nombre.IdDns)
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
			new_line := update_info.Value + " IN A " + antigua_ip
			txtlines[index] = new_line
		} else {
			new_line := update_info.Name + "." + update_info.Domain + " IN A " + update_info.Value
			txtlines[index] = new_line
		}

		// Se escribe el cambio en el log
		DomainLog(update_info.Name, update_info.Domain, update_info.Value, "update", update_info.IdDns)

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
		fmt.Println("Se creo con modifico! Reloj dominio " + update_info.Domain + ": " + reloj_a_string(ultimo_reloj))
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

		if n_lines != 0 {
			for _, eachline := range txtlines[:n_lines-1] {
				f.WriteString(eachline + "\n")
			}
			f.WriteString(txtlines[n_lines-1])
		}

		// Actualizar reloj
		s.Relojes[delete_info.Domain] = sumar_uno_a_reloj(s.Relojes[delete_info.Domain], int(delete_info.IdDns))

		// Enviar reloj como respuesta
		ultimo_reloj := s.Relojes[delete_info.Domain]
		fmt.Println("Se creo con elimino! Reloj dominio " + delete_info.Domain + ": " + reloj_a_string(ultimo_reloj))
		reloj_mensaje := &ClockMessage{X: int64(ultimo_reloj.X), Y: int64(ultimo_reloj.Y), Z: int64(ultimo_reloj.Z)}

		// Se escribe el cambio en el log
		DomainLog(delete_info.Name, delete_info.Domain, "", "delete", delete_info.IdDns)

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

func (s *Server) GetName(ctx context.Context, nombre *NewName) (*CommandResponse, error) {
	//-----------------------------------------------------------//
	//----------- EN ESTA PARTE SE BUSCA EL DNS -----------------//
	//----------- CON EL DOMINIO SOLICITADO Y SE RETORNA LA IP --//
	//----------- Y EL RELOJ VECTORIAL ASOCIADO -----------------//
	//-----------------------------------------------------------//

	file_name := zf_folder_paths[nombre.IdDns] + "/" + nombre.Domain + ".zf"
	// Chequear si el dominio existe. Esto es true si no existe
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		return &CommandResponse{Body: "Error, no existe el Dominio", Clock: nil}, nil
	} else {
		// Leer el archivo, leer linea por linea, y si el nombre no existe es tamos mal.}
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

		ip_addr := " "
		for _, eachline := range txtlines[:] {
			//fmt.Println(eachline)
			lname := strings.Split(strings.Split(eachline, " ")[0], ".")[0]
			if lname == nombre.Name {
				ip_addr = strings.Split(eachline, " ")[3]
				break
			}
		}
		if ip_addr != " " {
			dom := strings.Split(nombre.Domain, ".")[0]
			val := s.Relojes[dom]
			return &CommandResponse{Body: ip_addr, Clock: &ClockMessage{X: int64(val.X), Y: int64(val.Y), Z: int64(val.Z)}}, nil
		} else {
			return &CommandResponse{Body: "Error, no existe el Nombre", Clock: nil}, nil
		}
	}
}

func (s *Server) PropagarCambios(ctx context.Context, id_dns *IdDns) (*Message, error) {
	zf_path := zf_folder_paths[id_dns.IdDns]
	var log_files []string
	//En esta map se guardaran los dominios con una lista de los nombres modificados por dominio
	//var modificaciones map[string][]string
	//modificaciones = make(map[string][]string)

	fmt.Println("Propagandos cambios hacia el servidor 1...")

	err := filepath.Walk(zf_path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".log" {
			log_files = append(log_files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// leer cada archivo logs y ejecutar cada comando en el server 1
	for _, file := range log_files {
		file_name := file
		file, err := os.Open(file_name)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}

		fmt.Println("Leyendo " + file.Name() + " ...")

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var txtlines []string

		for scanner.Scan() {
			txtlines = append(txtlines, scanner.Text())
		}

		file.Close()

		// Conectarse al dns 1 para enviarle los cambios
		var conn_dns *grpc.ClientConn
		conn_dns, err = grpc.Dial(id_dns.IpDns, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn_dns.Close()

		s_dns1 := NewDnsServiceClient(conn_dns)

		for _, eachline := range txtlines[:] {
			fmt.Println("Enviando comando: " + eachline)
			params := strings.Split(eachline, " ")
			option := params[0]

			// Extraer nombre y dominio
			name := params[1]
			name_split := strings.Split(name, ".")

			if option == "create" {
				// Aqui comunicarse con el BROKER y obtener una ip de un dns
				new_ip := params[2]

				// Enviar al servidor dns el nombre que se quiere crear
				new_name := NewName{Name: name_split[0], Domain: name_split[1], Ip: new_ip, IdDns: int64(0)}

				// Pedir el vector reloj del dominio

				response, err := s_dns1.CreateName(context.Background(), &new_name)
				if err != nil {
					log.Fatalf("Error al tratar de crear nombre: %s", err)
				}
				log.Printf("Response from Server: %s", response.Body)
				// Si el reloj recibido no es nulo, se guarda como ultimo reloj de ese dominio
			}
			if option == "update" {
				var option string
				if len(strings.Split(params[2], ".")) == 4 {
					option = "ip"
				} else {
					option = "name"
				}
				update_info := UpdateInfo{Name: name_split[0], Domain: name_split[1], Opt: option, Value: params[2], IdDns: int64(0)}

				response, err := s_dns1.Update(context.Background(), &update_info)
				if err != nil {
					log.Fatalf("Error al tratar de updatear nombre: %s", err)
				}
				log.Printf("Response from Server: %s", response.Body)
			}
			if option == "delete" {
				delete_info := DeleteInfo{Name: name_split[0], Domain: name_split[1], IdDns: int64(0)}

				response, err := s_dns1.Delete(context.Background(), &delete_info)
				if err != nil {
					log.Fatalf("Error al tratar eliminar nombre: %s", err)
				}
				fmt.Printf("Response from Server: %s\n", response.Body)

			}
		}
	}
	fmt.Println("Cambios enviados al server 1 con exito.")
	return &Message{Body: "Cambios enviados al server 1 con exito."}, nil
}

func (s *Server) SobreescribirZf(ctx context.Context, zf_file *ZfFile) (*Message, error) {
	aux_split := strings.Split(zf_file.Nombre, "\\")
	name := aux_split[len(aux_split)-1]
	file_name := zf_folder_paths[zf_file.IdDns] + "/" + name
	f, err := os.Create(file_name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	f.WriteString(zf_file.Contenido)
	return &Message{Body: "ZFfile sobreescrita en server con exito."}, nil
}

// Funcion que solo es llamada por el servidor 1 envia los cambios al resto de servidores
func (s *Server) PropagarZfs(ctx context.Context, target_ips *TargetIps) (*Message, error) {

	zf_path := zf_folder_paths[target_ips.IdDns]
	var zf_files []string

	fmt.Println("Enviando ZFs al servidor 2...")

	err := filepath.Walk(zf_path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".zf" {
			zf_files = append(zf_files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	var conn_dns2 *grpc.ClientConn
	conn_dns2, err = grpc.Dial(target_ips.Ip1, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn_dns2.Close()

	s_dns2 := NewDnsServiceClient(conn_dns2)

	var conn_dns3 *grpc.ClientConn
	conn_dns3, err = grpc.Dial(target_ips.Ip2, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn_dns3.Close()

	s_dns3 := NewDnsServiceClient(conn_dns3)

	// leer cada archivo zf y ejecutar cada comando en el server 2 y 3
	for _, file := range zf_files {
		file_name := file
		file, err := os.Open(file_name)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}

		fmt.Println("Leyendo " + file.Name() + " ...")

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var txtlines string = ""

		for scanner.Scan() {
			if txtlines == "" {
				txtlines = txtlines + scanner.Text()
			} else {
				txtlines = txtlines + "\n" + scanner.Text()
			}

		}

		// enviar los cambios a los dns 2 y3
		log.Println("Enviando " + file_name + " al server 2...")
		var zf_file2 = ZfFile{IdDns: 1, Nombre: file_name, Contenido: txtlines}
		response, err := s_dns2.SobreescribirZf(context.Background(), &zf_file2)
		if err != nil {
			log.Fatalf("Error al tratar sobreescribir nombre: %s", err)
		}
		fmt.Printf("Response from Server: %s\n", response.Body)

		log.Println("Enviando " + file_name + " al server 3...")
		var zf_file3 = ZfFile{IdDns: 2, Nombre: file_name, Contenido: txtlines}
		response, err = s_dns3.SobreescribirZf(context.Background(), &zf_file3)
		if err != nil {
			log.Fatalf("Error al tratar sobreescribir nombre: %s", err)
		}
		fmt.Printf("Response from Server: %s\n", response.Body)

		file.Close()
	}

	// Leer todos los .zf
	return &Message{Body: "ZFfile sobreescrita en server con exito."}, nil
}

func (s *Server) EliminarLogs(ctx context.Context, id_dns *IdDns) (*Message, error) {
	zf_path := zf_folder_paths[id_dns.IdDns]
	var log_files []string
	//En esta map se guardaran los dominios con una lista de los nombres modificados por dominio
	//var modificaciones map[string][]string
	//modificaciones = make(map[string][]string)

	fmt.Println("Eliminando todos los logs.")

	err := filepath.Walk(zf_path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".log" {
			log_files = append(log_files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// leer cada archivo logs y ejecutar cada comando en el server 1
	for _, file := range log_files {
		os.Remove(file)
	}

	return &Message{Body: "Logs eliminados..."}, nil
}
