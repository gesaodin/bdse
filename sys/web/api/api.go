//interfaz de programación de aplicaciones.
//Abreviado como API del inglés: Application Programming Interface,
//es un conjunto de subrutinas, funciones y procedimientos
//(o métodos, en la programación orientada a objetos) que
//ofrece cierta biblioteca para ser utilizado por otro software
//como una capa de abstracción.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	balance "github.com/gesaodin/bdse/mdl/balance"
	"github.com/gesaodin/bdse/sys/seguridad"
)

//Usuario del sistema
type Usuario struct {
	Usr seguridad.Usuario `json:"usuario"`
}

//Pago para el control
type Pago struct{}

//Movimiento Ingreso y Egresos
type Movimiento struct{}

//Localizacion ubicacion geografica
type Localizacion struct{}

//Registro de Control
type Registro struct{}

//Comercializadora compuesta por grupo, subgrupo, colector, agencia
type Comercializadora struct{}

//Transferencia Solicitud de operaciones bancarias
type Transferencia struct{}

var pago balance.Pago

//Salvar un registro web
func (a *Usuario) Salvar(w http.ResponseWriter, r *http.Request) {

}

//ConsultarToken WEB
func (a *Usuario) ConsultarToken(w http.ResponseWriter, r *http.Request) {
}

//Salvar un pago
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

//GenerarCobrosYPagos De las consultas web
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
	//fmt.Println("Entrando...")
	j, e := pago.GenerarCobrosYPagos(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//GenerarCobrosYPagosGrupo De las consultas web
func (p *Pago) GenerarCobrosYPagosGrupo(w http.ResponseWriter, r *http.Request) {
	Cabecera(w, r.Header.Get("Origin"))
	var pagos balance.Pago
	err := json.NewDecoder(r.Body).Decode(&pagos)

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
	j, e := pagos.GenerarCobrosYPagosGrupo()
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//GenerarCobrosYPagosSistemas Programas MATICLOT, MORPHEUS, POS
func (p *Pago) GenerarCobrosYPagosSistemas(w http.ResponseWriter, r *http.Request) {
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
	j, e := pago.GenerarCobrosYPagosSistemas(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//GenerarCobrosYPagosDetallados de todos los programas
func (p *Pago) GenerarCobrosYPagosDetallados(w http.ResponseWriter, r *http.Request) {
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
	j, e := pago.GenerarCobrosYPagosDetallados(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//ListarPagos de los entregados y recibidos
func (p *Pago) ListarPagos(w http.ResponseWriter, r *http.Request) {
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
	j, e := pago.ListarPagos(dataJSON)
	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//CierreDiario Pos los saldos acumulados
func (p *Pago) CierreDiario(w http.ResponseWriter, r *http.Request) {
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
	j, e := pago.GenerarCierreDiario(dataJSON)

	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//CierreDiario Pos los saldos acumulados
func (p *Pago) CierreDiarioCalculo(w http.ResponseWriter, r *http.Request) {
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
	j, e := pago.GenerarCierreDiario(dataJSON)

	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//EstadoDeCuentaGrupo Pos los saldos acumulados
func (p *Pago) EstadoDeCuentaGrupo(w http.ResponseWriter, r *http.Request) {
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
	j, e := dataJSON.EstadoDeCuentasGrupo()

	if e != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error al consultar los datos"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
