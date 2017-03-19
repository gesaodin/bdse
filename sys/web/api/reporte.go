package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/mdl/loteria"
	"github.com/gesaodin/bdse/sys/seguridad"
)

type (
	Reporte struct {
	}
)

var ReporteLoteria loteria.Reporte

func Cabecera(w http.ResponseWriter, origen string) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", origen)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

}

func CabeceraRechazada(w http.ResponseWriter, estatus int, m string) {
	w.WriteHeader(estatus)
	msj := []byte(m)
	w.Write(msj)
}

func (a *Reporte) ReporteLoteriaArchivo(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var dataJSON loteria.JsonDataReporte
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
	j, e := ReporteLoteria.ArchivosCargados(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func (a *Reporte) ReporteSaldos(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))

	var dataJSON loteria.JsonDataReporte
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
	j, e := ReporteLoteria.Saldos(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func (a *Reporte) SaldosGeneralesPorSistema(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))

	var dataJSON loteria.JsonDataReporte
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
	j, e := ReporteLoteria.SaldosGeneralesPorSistemas(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

//SaldosGeneralesTotales Generales
func (a *Reporte) SaldosGeneralesTotales(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))

	var dataJSON loteria.JsonDataReporte
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
	j, e := ReporteLoteria.SaldoGeneralTotales(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

//BalanceGeneral Reglas de balance
func (a *Reporte) BalanceGeneral(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))

	var dataJSON loteria.JsonDataReporte
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
	j, e := ReporteLoteria.BalanceGeneral(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
