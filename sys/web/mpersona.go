package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gesaodin/bdse/mdl"
	"github.com/gesaodin/bdse/util"
	"github.com/gesaodin/bdse/util/logs/mensaje"
)

//PersonaGET Obtener una persona (find)
func PersonaGET(w http.ResponseWriter, r *http.Request) {
	var Prs mdl.Persona
	var m mensaje.MSJ

	fmt.Println("Validando")

	id := r.URL.Query().Get("id")
	fmt.Println(r.URL.Query())
	if id != "" {
		// origin := r.Header.Get("Origin")
		// fmt.Println("Accediendo al documento (", id, ") desde: ", origin)
		e := Prs.ConsultarMGO(id)
		Prs.ListarPostgreSQL()
		util.Error(e)
		if Prs.Cedula != "" {
			m.Estatus = true
			var o = mensaje.OBJRespuesta{Persona: Prs, MSJ: m}
			j, _ := json.MarshalIndent(o, "", " ")
			w.WriteHeader(http.StatusOK)
			w.Write(j)
		} else {
			CabeceraRechazada(w, http.StatusNotFound, "No se encontro el identificador")
		}
	} else {
		CabeceraRechazada(w, http.StatusUnauthorized, "No se encontro definición del objeto")
	}

}

//PersonaPOST crear una persona (insert)
func PersonaPOST(w http.ResponseWriter, r *http.Request) {
	var m mensaje.MSJ
	var p mdl.Persona

	Cabecera(w, r.Header.Get("Origin"))

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		CabeceraRechazada(w, http.StatusNotFound, err.Error())
	} else {
		if p.Cedula == "" {
			CabeceraRechazada(w, http.StatusNotFound, "Debe indicar un objeto")
			return
		}
		p.ID = bson.NewObjectId()
		p.FechaDeCreacion = time.Now()
		e := p.SalvarMGO()
		fmt.Println("Creando una persona ", p.ID)

		if e == nil {
			m.Estatus = true
			//m.Numero = p.Id
			m.Descripcion = "Objeto Creado"
			j, _ := json.Marshal(m)
			w.WriteHeader(http.StatusCreated)
			w.Write(j)
		} else {
			CabeceraRechazada(w, http.StatusNotFound, e.Error())
		}

	}
}

//PersonaPUT Agregar
func PersonaPUT(w http.ResponseWriter, r *http.Request) {
	var m mensaje.MSJ
	var p mdl.Persona
	var datos map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&datos)
	util.Error(err)
	fmt.Println("Entrando en el metodo PUT")

	datos["fechadecreacion"] = time.Now()
	e := p.ActualizarMGO(datos)
	if e == nil {
		fmt.Println("Actualizando una persona \n ", datos["cedula"])
		m.Estatus = true
		m.Descripcion = "Objeto Actualizado"
		j, _ := json.Marshal(m)
		w.WriteHeader(http.StatusCreated)
		w.Write(j)
	} else {
		CabeceraRechazada(w, http.StatusNotFound, e.Error())
	}

}

//PersonaUpdate Actualizar
func PersonaUpdate(w http.ResponseWriter, r *http.Request) {
	// var m logs.MSJ
	Cabecera(w, r.Header.Get("Origin"))
	w.WriteHeader(http.StatusCreated)
	fmt.Println("Entrando en el metodo PUT")
	/**
	if b == "bien" {
		//CabeceraAceptada(w, http.StatusCreated)
		fmt.Println("Entrando en el metodo PUT")
		// m.Estatus = true
		// m.Descripcion = "Objeto Actualizado"
		// j, _ := json.Marshal(m)

		//w.Write(j)
	} else {
		// m.Estatus = true
		// m.Descripcion = "Objeto No actualizado"
		//CabeceraAceptada(w, http.StatusCreated)
		// j, _ := json.Marshal(m)

		//w.Write(j)
	}
	*/
}

//PersonaDELETE Eliminar
func PersonaDELETE(w http.ResponseWriter, r *http.Request) {

}
