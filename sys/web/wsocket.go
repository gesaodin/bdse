package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gesaodin/bdse/util"
	"github.com/gesaodin/bdse/util/logs/mensaje"
	"github.com/gorilla/websocket"
)

type ()

var (
	Mensajeria = New()
	Upgrader   = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsocket = mensaje.WSocket{}
)

func EscribirMensajes(conn *websocket.Conn, usuario string, error bool) {
	//var wsocketCliente = mensaje.WSocket{}
	var m = mensaje.WChat{}
	var r = mensaje.WChat{}
	var entregado = mensaje.WChat{}
	for {

		e := conn.ReadJSON(&m)
		if e != nil {
			// fmt.Println(e)
			conn.Close()
			Mensajeria.EliminarUsuario(usuario)
			return
		}
		// fmt.Println(msj)

		r.Msj = m.Msj
		r.De = usuario
		r.Tiempo = time.Now()
		r.Tipo = 3

		entregado.Para = usuario
		entregado.Tiempo = time.Now()
		entregado.Tipo = 4

		if Mensajeria.Usuario[m.Para].ch != nil {
			j, _ := json.Marshal(r)
			Mensajeria.Usuario[m.Para].ch <- j
			entregado.Msj = "Entregado"
			jr, _ := json.Marshal(entregado)
			conn.WriteMessage(websocket.TextMessage, jr)
		} else {
			entregado.Msj = "No existe conexion con el destino"
			jr, _ := json.Marshal(entregado)
			conn.WriteMessage(websocket.TextMessage, jr)
		}
		// messageType, p, err := conn.ReadMessage()
		//util.Error(err)
		// if err != nil {
		// 	Mensajeria.EliminarUsuario(usuario)
		// 	conn.Close()
		// 	// conn.WriteMessage(messageType, []byte("Usuario ya existe"))
		// 	return
		// }
		// fmt.Println(p)
		// e := json.Unmarshal(p, msj)
		// if e != nil {
		// 	fmt.Println(e)
		// 	return
		// }
		// fmt.Println(msj.Msj)

		//wsocket.Listado = Mensajeria.ListarUsuariosMenos(usuario)
		// wsocketCliente.Usuario = Mensajeria.Usuario[usuario]
		// m := string(p)
		// wsocketCliente.Mensaje = m
		// j, _ := json.Marshal(wsocketCliente)
		// conn.WriteMessage(messageType, j)

		// if usuario == "carlos" {
		// 	Mensajeria.Usuario["yasmin"].ch <- j
		// }
		// fmt.Println("Listado de Usuario: ", mensajeria)
	}

}

func LogicaDelMensajePorTiempo(conn *websocket.Conn, usuario string) {
	ch := time.Tick(3 * time.Second)
	for range ch {
		conn.WriteMessage(websocket.TextMessage, []byte("Actualizando"))
	}
}

func LogicaDelMensajePorCanales(conn *websocket.Conn, usuario string, ch chan []byte) {

	for {
		select {
		case dato := <-ch:

			conn.WriteMessage(websocket.TextMessage, dato)
		}
	}
}

func CreandoWS(w http.ResponseWriter, r *http.Request) {
	usuario := r.URL.Query().Get("id")
	conn, err := Upgrader.Upgrade(w, r, nil)
	error := Mensajeria.CrearUsuario(usuario)
	util.Error(err)

	wsocket.Listado = Mensajeria.ListarUsuariosMenos(usuario)
	wsocket.Usuario = Mensajeria.Usuario[usuario]

	wsocket.Mensaje = "Se Establecio la conexión con el servidor"
	j, _ := json.Marshal(wsocket)
	conn.WriteMessage(websocket.TextMessage, j)

	// go LogicaDelMensajePorTiempo(conn, usuario)
	go LogicaDelMensajePorCanales(conn, usuario, Mensajeria.Usuario[usuario].ch)
	go EscribirMensajes(conn, usuario, error)

}

func MensajesWS() {

}
