package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/mdl/transferencia"
  "github.com/gesaodin/bdse/sys/seguridad"
)


//SalvarGrupo del sistema
func (t *Transferencia) Registrar(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var transf transferencia.Transferencia
	err := json.NewDecoder(r.Body).Decode(&transf)

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
	j, e := transf.Registrar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ListarAgencia del sistema
func (t *Transferencia) ListarAgencia(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var transf transferencia.Transferencia
	err := json.NewDecoder(r.Body).Decode(&transf)

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
	j, e := transf.ListarAgencia()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos: " + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}


//ListarGrupo del sistema
func (t *Transferencia) ListarGrupo(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var transf transferencia.Transferencia
	err := json.NewDecoder(r.Body).Decode(&transf)

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
	j, e := transf.ListarGrupo()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
