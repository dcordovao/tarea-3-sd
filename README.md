# sd-tarea3
Tarea sobre un sistema DNS implementando consistencia

## Integrantes

Sebastián Sanchez Lagos 201504022-2
Diego Córdova Opazo 201403009-6

## Instrucciones de ejecución
+ Primero se debe ejecutar en 3 máquinas distintas el comando make, luego en una 4 máquina ejecutar make broker

+ Una vez iniciados estos 4 componentes, elegir cualquiera de las maquinas de DNS, abrir otra terminal y ejecutar el comando make admin. Después, elegir cualquiera de las maquinas de DNS abrir otra terminal y ejecutar el comando make cliente.

+ Cada vez que se desee ejecutar todo el programa nuevamente, ejecutar make clean en los DNS y luego make para ejecutar nuevamente.

## El sistema consiste en:

+ Administradores que se encarga de crear nuevos registros ZF de dominio de los servidores DNS, además de agregar, actualizar o eliminar lı́neas en estos.

+ Servidores DNS (3 réplicas) que permiten almacenar registros ZF y llevar un Log de Cambios junto a cada registro ZF, cuya función será registrar los cambios (create/delete/update) que se haya realizado el Administrador en el registro ZF.

+ Un Broker encargado de balancear la carga entre los diversas réplicas. Actúa como un intermediario entre los Administradores y Clientes y los Servidores DNS.

+ Clientes que realizan las consultas al Broker para saber la dirección IP de la página web que está solicitando.

## Funcionamiento a grandes rasgos

##Supuestos y explicaciones importantes en la tarea:

+ DNS
+ Explicación: 
	+ Create nombre
		+ Si no existe el dominio para ese nombre crea los archivos <dominio.zf> y <dominio.log>, se añade el nombre al <dominio.zf> y se agrega el comando al <dominio.log>.
		+ Si ya existe el dominio pero no el nombre, se añade el nombre al <dominio.zf> y se agrega el comando al <dominio.log>.
		+ Si ya existe el dominio y el nombre, no se agrega ni se crea nada y se retorna un mensaje de alerta auto descriptivo.

	+ Update nombre
		+ Si el nombre y el dominio existen, cada linea del archivo <dominio.zf> respectivo se guarda en una lista, y la posición del nombre en el archivo se guarda en una variable. Luego el archivo se borra y se vuelve a crear pero esta vez actualizando el nombre en la posición previamente guardada. Luego se agrega la operacion al log.

	+ Delete nombre
		+ Si el nombre y el dominio existen, cada linea del archivo <dominio.zf> respectivo se guarda en una lista, excepto el elemento que se va a borrar. Luego el archivo se borra y se vuelve a crear pero esta vez el elemento ya no esta en la lista, por lo tanto su linea en el archivo ya no existe. Luego se agrega la operacion al log.


+ Supuestos:


+ Broker
+ Explicación:
	+ Redirige la conexion entre el Cliente y los DNS, en primera instancia de froma aleatoria, después de forma determinada (ver explicación del cliente)
	+ Entrega ip de DNS al azar al administrador. 

+ Supuestos: 


+ Administrador(es)
+ Explicación: 
	+ escribir en la terminal cualquiera de estos comandos:
		+ create <nombre.dominio> <ip>, ej.: create cristal.cl 10.9.8.7
		+ update <nombre.dominio> ip <ip>, ej.: update cristal.cl ip 10.9.9.8
		+ update <nombre1.dominio> name <nombre2.dominio>, ej.: update cristal.cl name escudo.cl
		+ delete <nombre.dominio>, ej.: delete deepl.com
	+ Se ingresa y valida el comando escogido.
	+ Se envia la solicitud al broker
	+ El broker retorna una ip al azar y luego el admin se conecta directamente al dns escogido.
	+ Se verifica que los relojes calcen con el ultimo leido guardado en memoria, si no se cumple, o el ip otorgado por el broker no coincide con la instrucción que se desea ejecutar, entonces se utiliza la ip asociada al dominio que se quiere crear, actualizar o borrar. 

+ Supuestos:	


+ Cliente(s)
	+ Explicación: 			
		+ Se ingresa y valida el comando get <nombre.dominio>.
		+ Se envía la solicitud al broker.
		+ El broker retorna la respuesta del dns al azar de vuelta al cliente.
		+ Si se encontró lo solicitado, el cliente imprime ip y vector, luego compara dicho vector con el que se halla su diccionario de dominio:(reloj, ip) para ver que sea el ultimo, si esta desactualizado envia la solicitud nuevamente al broker, pero con la ip del dicionario.
		+ Si no lo encontró, envia la solicitud nuevamente al broker pero con la ip del dicionario el cliente. Luego si no esta, significa que fue recientemente borrado, de lo contrario imprime ip y actualiza vector.

	+ Supuestos: solo se recuerda el reloj de las lecturas que tuvieron exito, es decir en las cuales se encontró el nombre.dominio requerido.



## Estructura del proyecto

```bash
.
├── admin
│   └── admin.go
├── broker
│   └── broker_server.go
├── cliente
│   └── cliente.go
├── go.mod
├── go.sum
├── grpc_services
│   ├── broker_service
│   │   ├── broker_service.go
│   │   └── broker_service.pb.go
│   ├── broker_service.proto
│   ├── dns_service
│   │   ├── dns_service.go
│   │   └── dns_service.pb.go
│   ├── dns_service.proto
│   └── Makefile
├── README.md
└── servidor_dns
    ├── servidor_dns.go
    ├── zf_files1
    │   └── dummy_file
    ├── zf_files2
    │   └── dummy_file
    └── zf_files3
        └── dummy_file
```





## Coneccion con Red DI y Ejecución en Máquinas virtuales

IMPORTANTE: Consideraciones acerca del orden de ejecucion de las máquinas


+ Máquina 1: 
	+ Entidad:    DNS (y cliente(s) y/o admin(s))
	+ ip:         10.10.28.121
	+ contraseña: DSmvWkQsyIkaJzU


+ Máquina 2:
	+ Entidad:    DNS (y cliente(s) y/o admin(s))
	+ ip:         10.10.28.122
	+ contraseña: eDtGthpFSmaypHj


+ Máquina 3:
	+ Entidad:    DNS (y cliente(s) y/o admin(s))
	+ ip:         10.10.28.123
	+ contraseña: kndFkwYEQRdcTTu


+ Máquina 4:
	+ Entidad:    Broker
	+ ip:         10.10.28.124
	+ contraseña: XvfDKuTBWbAXEgj

Se creo un go mod en "github.com/dcordova/sd_tarea3" y se puede acceder a los servicios:
	+ "github.com/dcordova/sd_tarea3/broker_service"
	+ "github.com/dcordova/sd_tarea3/dns_service"

Perdon por no poner tildes en  los print, pero no usamos porque se podia caer la wea por caracteres especiales

