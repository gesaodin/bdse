package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	balance "github.com/gesaodin/bdse/mdl/balance"
	"github.com/gesaodin/bdse/sys/seguridad"
)

type Usuario struct {
	Usr seguridad.Usuario `json:"usuario"`
}

type Pago struct{}

type Movimiento struct{}

var pago balance.Pago

func (a *Usuario) Salvar(w http.ResponseWriter, r *http.Request) {

}

func (a *Usuario) ConsultarToken(w http.ResponseWriter, r *http.Request) {

}

func (p *Pago) Salvar(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var dataJSON balance.Pago
	err := json.NewDecoder(r.Body).Decode(&dataJSON)

	if err != nil {
		fmt.Println("Estoy en un error ", err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}

	_, e := seguridad.Stores.Get(r, "session-bdse")
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error en la Cookies"))
		return
	}
	j, e := pago.Registrar(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (p *Pago) GenerarCobrosYPagos(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var dataJSON balance.Pago
	err := json.NewDecoder(r.Body).Decode(&dataJSON)

	if err != nil {
		fmt.Println("Estoy en un error ", err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}

	_, e := seguridad.Stores.Get(r, "session-bdse")
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error en la Cookies"))
		return
	}
	j, e := pago.GenerarCobrosYPagos(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
