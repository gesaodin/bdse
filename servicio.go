/*
Copyright 2017 Carlos Peña.Todos los derechos reservados.

En informática un Bus de Servicio Empresarial (ESB por sus siglas en inglés)
es un modelo de arquitectura de software que gestiona la comunicación entre
servicios web. Es un componente fundamental de la Arquitectura Orientada a
Servicios.

Un ESB generalmente proporciona una capa de abstracción construida sobre
una implementación de un sistema de mensajes de empresa que permita a los
expertos en integración explotar el valor del envío de mensajes sin tener que
escribir código. Al contrario que sucede con la clásica integración de
aplicaciones de empresa (IAE) que se basa en una pila monolítica sobre una
arquitectura hub and spoke, un bus de servicio de empresa se construye sobre
unas funciones base que se dividen en sus partes constituyentes, con una
implantación distribuida cuando se hace necesario, de modo que trabajen
armoniosamente según la demanda.

*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/sys/seguridad"
	"github.com/gesaodin/bdse/sys/web"
	"github.com/gorilla/context"
)

func init() {
	fmt.Println("")
	fmt.Println("Versión del Panel ", sys.Version)
	fmt.Println("")
	if sys.MongoDB {
		fmt.Println("Metodo de Encriptamiento ", seguridad.Encriptamiento, "...")
		// sys.MongoDBConexion()
		sys.PostgresDBConexion()
		fmt.Println("")
		fmt.Println("..........................................................")
		fmt.Println("... Iniciando Carga de Elemento Para el servidor WEB   ...")
		fmt.Println("..........................................................")
		fmt.Println("")
	}
}

func main() {
	// var archivo = util.Archivo{}

	//MORPHEUS LOTERIA
	// archivo.Ruta = "public/test/loteria/"
	// archivo.NombreDelArchivo = "Mo.Txt"
	// archivo.LeerMorpheus(sys.PostgreSQL)
	// res, err := sys.PostgreSQL.Exec(s)
	// fmt.Println(res, err)

	//POS (1,2,3) LOTERIA
	// archivo.NombreDelArchivo = "P1.txt"
	// _, s := archivo.LeerPos()
	// res, err := sys.PostgreSQL.Exec(s)
	// fmt.Println(res, err)

	//MATICLO LOTERIA
	// archivo.NombreDelArchivo = "Ma.xlsx"
	// _, s := archivo.LeerMaticlo()
	// res, err := sys.PostgreSQL.Exec(s)
	// fmt.Println(res, err)

	//ILBANQUERO PARLEY
	// archivo.Ruta = "public/test/parley/"
	// archivo.NombreDelArchivo = "Il.csv"
	// _, s := archivo.LeerIlbanquero()
	// res, err := sys.PostgreSQL.Exec(s)
	// fmt.Println(res, err)

	//SPORT PARLEY
	// archivo.Ruta = "public/test/parley/"
	// archivo.NombreDelArchivo = "Sp.txt"
	// _, s := archivo.LeerSport()
	// res, err := sys.PostgreSQL.Exec(s)
	// fmt.Println(res, err)

	// fmt.Println("Leyendo Archivos")

	web.CargarModulosWeb()
	go http.ListenAndServeTLS(":3000", "sys/seguridad/https/cert.pem", "sys/seguridad/https/key.pem", web.WsEnrutador)
	fmt.Println("Servidor Escuchando en el puerto:  3000")

	//http://dominigy o.com/*
	srv := &http.Server{
		Handler:      context.ClearHandler(web.Enrutador),
		Addr:         ":" + sys.Puerto,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Servidor Escuchando en el puerto: ", sys.Puerto)
	go srv.ListenAndServe()

	//https://dominio.com/* Protocolo de capa de seguridad
	server := &http.Server{
		Handler:      context.ClearHandler(web.Enrutador),
		Addr:         ":" + sys.PuertoSSL,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Servidor Escuchando en el puerto: ", sys.PuertoSSL)
	log.Fatal(server.ListenAndServeTLS("sys/seguridad/https/cert.pem", "sys/seguridad/https/key.pem"))

}
