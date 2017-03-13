package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/mdl/localizacion"
	"github.com/gesaodin/bdse/sys/seguridad"
)

//ConsultarEstado Localizacion de Estados
func (l *Localizacion) ConsultarEstado(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var localizacion localizacion.Estado
	err := json.NewDecoder(r.Body).Decode(&localizacion)

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
	j, e := localizacion.Consultar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ConsultarCiudad Localizacion de Estados
func (l *Localizacion) ConsultarCiudad(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var localizacion localizacion.Ciudad
	err := json.NewDecoder(r.Body).Decode(&localizacion)

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
	j, e := localizacion.Consultar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ConsultarMunicipio Localizacion de Estados
func (l *Localizacion) ConsultarMunicipio(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var localizacion localizacion.Municipio
	err := json.NewDecoder(r.Body).Decode(&localizacion)

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
	j, e := localizacion.Consultar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ConsultarParroquia Localizacion de Estados
func (l *Localizacion) ConsultarParroquia(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var localizacion localizacion.Parroquia
	err := json.NewDecoder(r.Body).Decode(&localizacion)

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
	j, e := localizacion.Consultar()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
