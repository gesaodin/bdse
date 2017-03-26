package web

//Copyright Carlos Peña
//Modulo de negociación WEB
import (
	"fmt"
	"net/http"

	"github.com/gesaodin/bdse/sys/web/api"
	"github.com/gorilla/mux"
)

//Variables de Control
var (
	Enrutador   = mux.NewRouter()
	WsEnrutador = mux.NewRouter()
)

//CargarModulosWeb Cargador de modulos web
func CargarModulosWeb() {
	WChat()
	WMSeguridad()
	WMPersona()
	WAPI()
	WMAdminLTE()
}

//Principal Página inicial del sistema o bienvenida
func Principal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Saludos bienvenidos al Bus Empresarial de Datos")
}

//WMSeguridad Esquema de Seguridad Web
func WMSeguridad() {
	Enrutador.HandleFunc("/", Principal)
	Enrutador.HandleFunc("/admin/login", Login)
	// Enrutador.HandleFunc("/admin/loginToken", LoginToken)
	// Enrutador.HandleFunc("/admin/session", seguridad.Session)
	Enrutador.HandleFunc("/admin/logout", Login)
	fmt.Println("Cargando Modulos de Seguridad (Session, Token, Cookies)...")
}

//WMAdminLTE OpenSource tema de panel de control
//Tecnología Bootstrap3
func WMAdminLTE() {
	fmt.Println("Cargando Modulos de AdminLTE...")
	var GP = GPanel{}
	Enrutador.HandleFunc("/bdse-admin/gpanel/{id}", GP.Login)
	prefix := http.StripPrefix("/bdse-admin/gpanel", http.FileServer(http.Dir("public_web/adminlte")))
	Enrutador.PathPrefix("/bdse-admin/gpanel/").Handler(prefix)
	prefixx := http.StripPrefix("/bdse-admin/public/temp", http.FileServer(http.Dir("public/temp")))
	Enrutador.PathPrefix("/bdse-admin/public/temp/").Handler(prefixx)
}

//WAPI Aplicaciones de interfaz de programacion
func WAPI() {
	var API = api.Reporte{}
	var Listar = api.Listar{}
	var Pago = api.Pago{}
	var Movimiento = api.Movimiento{}
	var Localizacion = api.Localizacion{}
	var Registro = api.Registro{}
	var Comercializadora = api.Comercializadora{}
	var Transferencia = api.Transferencia{}

	url := "/bdse-admin/gpanel/api/"

	Enrutador.HandleFunc(url+"reportearchivo", API.ReporteLoteriaArchivo)
	Enrutador.HandleFunc(url+"reportesaldo", API.ReporteSaldos).Methods("POST")
	Enrutador.HandleFunc(url+"reportesaldogeneral", API.SaldosGeneralesPorSistema).Methods("POST")
	Enrutador.HandleFunc(url+"reportesaldototales", API.SaldosGeneralesTotales).Methods("POST")
	Enrutador.HandleFunc(url+"balancegeneral", API.BalanceGeneral).Methods("POST")

	Enrutador.HandleFunc(url+"balance/registrarpago", Pago.Salvar).Methods("POST")
	Enrutador.HandleFunc(url+"balance/cobrosypagos", Pago.GenerarCobrosYPagos).Methods("POST")
	Enrutador.HandleFunc(url+"balance/cobrosypagosgrupo", Pago.GenerarCobrosYPagosGrupo).Methods("POST")
	Enrutador.HandleFunc(url+"balance/cobrosypagossistemas", Pago.GenerarCobrosYPagosSistemas).Methods("POST")
	Enrutador.HandleFunc(url+"balance/cobrosypagosdetallados", Pago.GenerarCobrosYPagosDetallados).Methods("POST")
	Enrutador.HandleFunc(url+"balance/listarpagos", Pago.ListarPagos).Methods("POST")
	Enrutador.HandleFunc(url+"balance/cierrediario", Pago.CierreDiario).Methods("POST")
	Enrutador.HandleFunc(url+"balance/estadocuentagrupo", Pago.EstadoDeCuentaGrupo).Methods("POST")

	// API DE LISTADOS
	Enrutador.HandleFunc(url+"listasistema", Listar.Sistemas).Methods("POST")
	Enrutador.HandleFunc(url+"listasaldos", Listar.SaldosGeneral).Methods("POST")

	// MOVIMIENTOS
	Enrutador.HandleFunc(url+"movimiento/registrar", Movimiento.Registrar).Methods("POST")
	Enrutador.HandleFunc(url+"movimiento/listar", Movimiento.Listar).Methods("POST")
	Enrutador.HandleFunc(url+"movimiento/listardeposito", Movimiento.ListarDeposito).Methods("POST")
	Enrutador.HandleFunc(url+"movimiento/listarcuentas", Movimiento.ListarCuentas).Methods("POST")
	//Enrutador.HandleFunc(url+"movimiento/listarbanco", Movimiento.ListarBancos).Methods("GET")
	Enrutador.HandleFunc(url+"movimiento/actualizarer", Movimiento.ActualizarER).Methods("POST")

	//LOCALIZACION
	Enrutador.HandleFunc(url+"localizacion/consultarestado", Localizacion.ConsultarEstado).Methods("POST")
	Enrutador.HandleFunc(url+"localizacion/consultarciudad", Localizacion.ConsultarCiudad).Methods("POST")
	Enrutador.HandleFunc(url+"localizacion/consultarmunicipio", Localizacion.ConsultarMunicipio).Methods("POST")
	Enrutador.HandleFunc(url+"localizacion/consultarparroquia", Localizacion.ConsultarParroquia).Methods("POST")

	//GRUPO
	Enrutador.HandleFunc(url+"registro/grupo", Registro.SalvarGrupo).Methods("POST")
	Enrutador.HandleFunc(url+"registro/subgrupo", Registro.SalvarSubGrupo).Methods("POST")
	Enrutador.HandleFunc(url+"registro/colector", Registro.SalvarColector).Methods("POST")
	Enrutador.HandleFunc(url+"registro/agencia", Registro.SalvarAgencia).Methods("POST")

	//PERFIL DE LA Comercializadora
	Enrutador.HandleFunc(url+"perfil/comercializadora", Comercializadora.Consultar).Methods("POST")


	//TRANSFERENCIAS REGISTRO Y SOLICITUD
	Enrutador.HandleFunc(url+"transferencia/registrar", Transferencia.Registrar).Methods("POST")
	Enrutador.HandleFunc(url+"transferencia/listaagencia", Transferencia.ListarAgencia).Methods("POST")
	Enrutador.HandleFunc(url+"transferencia/listagrupo", Transferencia.ListarGrupo).Methods("POST")
}

//WChat Chat de Presentacion
func WChat() {
	WsEnrutador.HandleFunc("/", Principal)
	WsEnrutador.HandleFunc("/wsapi/c1", CreandoWS)
}
