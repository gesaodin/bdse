package web

import (
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/sys/web/api"
	"github.com/gorilla/mux"
)

var (
	Enrutador   = mux.NewRouter()
	WsEnrutador = mux.NewRouter()
)

func CargarModulosWeb() {
	WChat()
	WMSeguridad()
	WMPersona()
	WAPI()
	WMAdminLTE()
}

func Principal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Saludos bienvenidos al Bus Empresarial de Datos")
}

func WMSeguridad() {

	Enrutador.HandleFunc("/", Principal)
	Enrutador.HandleFunc("/admin/login", Login)
	// Enrutador.HandleFunc("/admin/loginToken", LoginToken)
	// Enrutador.HandleFunc("/admin/session", seguridad.Session)
	Enrutador.HandleFunc("/admin/logout", Login)

	fmt.Println("Cargando Modulos de Seguridad (Session, Token, Cookies)...")
}

func WMArchivosAngular() {
	fmt.Println("Cargando Modulos de Angular...")
}

func WMAdminLTE() {
	fmt.Println("Cargando Modulos de AdminLTE...")
	var GP = GPanel{}

	Enrutador.HandleFunc("/bdse-admin/gpanel/{id}", GP.Login)
	//Enrutador.HandleFunc("/bdse-admin/gpanel/subir", GP.SubirArchivoLoteria)
	prefix := http.StripPrefix("/bdse-admin/gpanel", http.FileServer(http.Dir("public_web/AdminLTE")))
	Enrutador.PathPrefix("/bdse-admin/gpanel/").Handler(prefix)

	prefixx := http.StripPrefix("/bdse-admin/public/temp", http.FileServer(http.Dir("public/temp")))
	Enrutador.PathPrefix("/bdse-admin/public/temp/").Handler(prefixx)

}

func WAPI() {
	var API = api.Reporte{}
	var Listar = api.Listar{}
	var Pago = api.Pago{}
	var Movimiento = api.Movimiento{}

	var base_api string = "/bdse-admin/gpanel/api/"
	Enrutador.HandleFunc(base_api+"reportearchivo", API.ReporteLoteriaArchivo)
	Enrutador.HandleFunc(base_api+"reportesaldo", API.ReporteSaldos).Methods("POST")
	Enrutador.HandleFunc(base_api+"reportesaldogeneral", API.SaldosGeneralesPorSistema).Methods("POST")
	Enrutador.HandleFunc(base_api+"reportesaldototales", API.SaldosGeneralesTotales).Methods("POST")
	Enrutador.HandleFunc(base_api+"balancegeneral", API.BalanceGeneral).Methods("POST")

	Enrutador.HandleFunc(base_api+"balance/registrarpago", Pago.Salvar).Methods("POST")
	Enrutador.HandleFunc(base_api+"balance/cobrosypagos", Pago.GenerarCobrosYPagos).Methods("POST")
	Enrutador.HandleFunc(base_api+"balance/listarpagos", Pago.ListarPagos).Methods("POST")

	// API DE LISTADOS
	Enrutador.HandleFunc(base_api+"listasistema", Listar.Sistemas).Methods("POST")
	Enrutador.HandleFunc(base_api+"listasaldos", Listar.SaldosGeneral).Methods("POST")

	// MOVIMIENTOS
	Enrutador.HandleFunc(base_api+"movimiento/registrar", Movimiento.Registrar).Methods("POST")
	Enrutador.HandleFunc(base_api+"movimiento/listardeposito", Movimiento.ListarDeposito).Methods("POST")

}

func WChat() {
	WsEnrutador.HandleFunc("/", Principal)
	WsEnrutador.HandleFunc("/wsapi/c1", CreandoWS)

}
