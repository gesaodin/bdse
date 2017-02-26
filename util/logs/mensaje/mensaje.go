package mensaje

import "time"

type (
	//Sistema
	MSJ struct {
		Estatus     bool   `json:"estatus"` //true: Exitoso | false: Ocurrio un Error
		Tipo        string `json:"tipo" `
		Numero      string `json:"numero"`
		Descripcion string `json:"descripcion"`
	}

	OBJRespuesta struct {
		MSJ     `json:"msj"`
		Persona interface{} `json:"persona"`
	}

	OBJMsj struct {
		MSJ `json:"msj"`
	}

	//Socket de Conexion por chat y sistema
	WSocket struct {
		Listado interface{} `json:"lst,omitempty"`
		Usuario interface{} `json:"usu,omitempty"`
		Mensaje string      `json:"msj"`
		Tipo    int         `json:"tipo"`
	}

	//Para el Chat
	WChat struct {
		Tipo   int       `json:"tipo,omitempty"`
		De     string    `json:"de,omitempty"`
		Para   string    `json:"para,omitempty"`
		Msj    string    `json:"msj"`
		Tiempo time.Time `json:"tiempo"`
	}
)
