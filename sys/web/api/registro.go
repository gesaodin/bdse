package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/mdl/grupo"
	"github.com/gesaodin/bdse/sys/seguridad"
)

//SalvarGrupo del sistema
func (re *Registro) SalvarGrupo(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var grupo grupo.Grupo
	err := json.NewDecoder(r.Body).Decode(&grupo)

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
	j, e := grupo.Registrar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//SalvarSubGrupo del sistema
func (re *Registro) SalvarSubGrupo(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var grupo grupo.Grupo
	err := json.NewDecoder(r.Body).Decode(&grupo)

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
	j, e := grupo.Registrar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//SalvarColector del sistema
func (re *Registro) SalvarColector(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var grupo grupo.Grupo
	err := json.NewDecoder(r.Body).Decode(&grupo)

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
	j, e := grupo.Registrar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//SalvarAgencia del sistema
func (re *Registro) SalvarAgencia(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var grupo grupo.Grupo
	err := json.NewDecoder(r.Body).Decode(&grupo)

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
	j, e := grupo.Registrar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
