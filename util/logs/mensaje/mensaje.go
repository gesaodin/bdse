//cadena de contenido entre servidor y clientes
package mensaje

import "time"

type (
	//MSJ mensajes
	MSJ struct {
		Estatus     bool   `json:"estatus"` //true: Exitoso | false: Ocurrio un Error
		Tipo        string `json:"tipo" `
		Numero      string `json:"numero"`
		Descripcion string `json:"descripcion"`
	}

	//OBJRespuesta respuestas
	OBJRespuesta struct {
		MSJ     `json:"msj"`
		Persona interface{} `json:"persona"`
	}

	//OBJMsj mensajes
	OBJMsj struct {
		MSJ `json:"msj"`
	}

	//WSocket conexion para chat y sistema
	WSocket struct {
		Listado interface{} `json:"lst,omitempty"`
		Usuario interface{} `json:"usu,omitempty"`
		Mensaje string      `json:"msj"`
		Tipo    int         `json:"tipo"`
	}

	//WChat el Chat
	WChat struct {
		Tipo   int       `json:"tipo,omitempty"`
		De     string    `json:"de,omitempty"`
		Para   string    `json:"para,omitempty"`
		Msj    string    `json:"msj"`
		Tiempo time.Time `json:"tiempo"`
	}
)
