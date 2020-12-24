# sd-tarea3
Tarea sobre consistencia implementando un sistema DNS

## Integrantes

Sebastián Sanchez Lagos 201504022-2
Diego Córdova Opazo 201403009-6

## Instrucciones ejecución

Definir el funcionamiento general del sistema distribuido final
IMPORTANTE: Consideraciones importantes en la tarea


## Estructura del proyecto

```bash
.
├── cliente
├────── cliente.go
```

OBS: Makefile correra cada nodo dependiendo de en que maquina nos encontremos.



## Funcionamiento a grandes rasgos


## Detalles y Consideraciones

## Coneccion con Red DI y Ejecución en Máquinas virtuales

IMPORTANTE: Consideraciones acerca del orden de ejecucion de las máquinas

* Para ejecutar cada maquina, hacer "cd tarea-3-sd" para ir al git luego solo correr el comando "make" (Esto correra un datanode o el namenode dependiendo de la máquina)
* Para correr un cliente downloader o uploader se debe hacer desde la máquina 1, debido a las ip's de las conexiones
* Para conectarse con un uploader o downloader desde otra maquina habría que cambiar las ip's en el código
de los archivos client.go correspondientes




+ Máquina 1: 
	+ ip:         10.10.28.121
	+ contraseña: DSmvWkQsyIkaJzU


+ Máquina 2:
	+ ip:         10.10.28.122
	+ contraseña: eDtGthpFSmaypHj


+ Máquina 3: 
	+ ip:         10.10.28.123
	+ contraseña: kndFkwYEQRdcTTu


+ Máquina 4: 
	+ ip:         10.10.28.124
	+ contraseña: XvfDKuTBWbAXEgj


# Usar proto
## Revisar esto a medida que se vaya avanzando!!

se creo un go mod en "github.com/dcordova/sd_tarea3" y se puede acceder a los servicios 
"github.com/dcordova/sd_tarea3/broker_service" y "github.com/dcordova/sd_tarea3/dns_service"

En caso de actualizar los .proto debemos hacer:
- protoc --go_out=plugins=grpc:dns_service dns_service.proto
- protoc --go_out=plugins=grpc:broker_service broker_service.proto
respectivamente, para actualizar los .pb.go

