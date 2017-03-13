package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/mdl/comercializadora"
	"github.com/gesaodin/bdse/sys/seguridad"
)

//DatosPerfil Esquema de datos web para el perfil
type DatosPerfil struct {
	Gastos              float64 `json:"gastos,omitempty"`
	CantidadAgencia     int     `json:"cant_agencia,omitempty"`
	CantidadGrupo       int     `json:"cant_grupo,omitempty"`
	CantidadSubGrupo    int     `json:"cant_subgrupo,omitempty"`
	DepositosPendientes float64 `json:"depositos,omitempty"`
}

//Consultar Localizacion de Estados
func (c *Comercializadora) Consultar(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var data DatosPerfil
	var comercializadora comercializadora.Comercializadora
	err := json.NewDecoder(r.Body).Decode(&comercializadora)

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
	m, _ := comercializadora.Cantidad()
	data.CantidadAgencia = m.Cantidad
	x, _ := comercializadora.Gastos()
	data.Gastos = x.Monto

	jSon, e := json.Marshal(data)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos" + e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jSon)

}
