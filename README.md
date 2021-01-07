# sd-tarea3
Tarea sobre un sistema DNS implementando consistencia

## Integrantes

Sebastián Sanchez Lagos 201504022-2
Diego Córdova Opazo 201403009-6

## Instrucciones de ejecución
+ Primero se debe ejecutar en 3 máquinas distintas el comando make firewall, luego poner la clave de esa maquina y despues ejecutar make.
+ Luego en una 4 máquina ejecutar make firewall, poner la clave y después make broker

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
	+ Create
		+ Si no existe el dominio para ese nombre crea los archivos <dominio.zf> y <dominio.log>, se añade el nombre al <dominio.zf> y se agrega el comando al <dominio.log>.
		+ Si ya existe el dominio pero no el nombre, se añade el nombre al <dominio.zf> y se agrega el comando al <dominio.log>.
		+ Si ya existe el dominio y el nombre, no se agrega ni se crea nada y se retorna un mensaje de alerta auto descriptivo.

	+ Update
		+ Si el nombre y el dominio existen, cada linea del archivo <dominio.zf> respectivo se guarda en una lista, y la posición del nombre en el archivo se guarda en una variable. Luego el archivo se borra y se vuelve a crear pero esta vez actualizando el nombre en la posición previamente guardada. Luego se agrega la operacion al log.

	+ Delete
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

## Monotonic Read

El monotonic reads consiste en el cliente al obtener un archivo de dominio guarda el reloj y el dominio en un diccionario de forma {dominio: relo, ip}, si nuevamente vuelve a obtener información de dicho dominio lo compara con el valor guardado en el diccionario. La comparación se realiza con la función compare_clocks(new_clock, last_clock). En caso de estar desactualizado se llama nuevamante a la funcion Connect pero ahora con una ip especfica, es decir el broker ya no tirara un numero al azar para conectarse, si no que tratara directamente de conectarse a dicha ip.

## Read Your Writes

Para implementarlo se tienen relojes tanto en los dns como en el admin, para cada dominio. En caso de que el admin modifique un modinio, este guarda en un diccionario el reloj y la ip del dominio midificado (las llaves del diccionario son dominios). Luego cada vez que se realiza un comando, el dns envia un reloj de vuelta, y el admin revisa si es que ya tiene un relol guardado para dicho dominio. En caso de ser verdad esto, el admin los compara. Si el vector que tiene guardado es mas actual (el reloj nuevo tiene una posicion x,y o z desactualizada) se cambia la conexion que se iba a realizar por una con la ip del reloj que se tiene guardado.

## Propagación

La propagación se implemento de manera que hay un servidor principal (el servidor 1 de ip 10.10.28.121), y el resto envia sus cambios a este. Estos cambios provienen de todos los archivos logs. El servidor principal realiza estos cambios en sus archivos, pero solo los cambios que no generan conflictos. Luego de esto, el servidor principal envia todos sus archivos ZF, que contienen las versiones mergeadas de todos cambios, a los servidores 2 y 3, y estos simplemente sobreescriben sus archivos con estos. 

El broker es el encargado de llamar cada 5 minutos a la función PropagarCambios de el servidor dns 1. Luego, esto activa todo lo explicado anteriormente. 

Comando propagate: el admin puede realizar el comando "propagate" para realizar los cambios sin tener que esperar 5 minutos.

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
    ├── zf_files
    │   └── dummy_file

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

