package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/mdl/movimiento"
	"github.com/gesaodin/bdse/sys/seguridad"
)

//Registrar movimientos
func (m *Movimiento) Registrar(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var movimiento movimiento.Movimiento
	err := json.NewDecoder(r.Body).Decode(&movimiento)

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
	j, e := movimiento.Salvar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ListarDeposito Bancarios
func (m *Movimiento) ListarDeposito(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var movimiento movimiento.Movimiento
	err := json.NewDecoder(r.Body).Decode(&movimiento)
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
	j, e := movimiento.ListarDepositos()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ListarCuentas Movimiento | Cuentas Bancarias | Ambos
func (m *Movimiento) ListarCuentas(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var movimiento movimiento.Movimiento
	err := json.NewDecoder(r.Body).Decode(&movimiento)
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
	j, e := movimiento.ListarCuentas()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//Listar Cuentas para Bancos
func (m *Movimiento) Listar(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var movimiento movimiento.Movimiento
	err := json.NewDecoder(r.Body).Decode(&movimiento)
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
	j, e := movimiento.Listar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ActualizarER Entregado Recibidos
func (m *Movimiento) ActualizarER(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var movimiento movimiento.Movimiento
	err := json.NewDecoder(r.Body).Decode(&movimiento)
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
	j, e := movimiento.Actualizar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
