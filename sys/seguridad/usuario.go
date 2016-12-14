package seguridad

import "time"

const (
	Administrador        = "Administrador"
	AdministradorDeGrupo = "AdministradorDeGrupo"
	Invitado             = "Invitado"
	Produccion           = "Produccion"
	Desarrollador        = "Desarrollador"
	Pasante              = "Pasante"
	Consulta             = "Consulta"
	Root                 = "Root" //Todos los privilegios del sistema
)

type MetodoSeguro struct {
	Consultar  bool `json:"consultar"`
	Insertar   bool `json:"insertar"`
	Actualizar bool `json:"actualizar"`
	Eliminar   bool `json:"eliminar"`
	Crud       bool `json:"crud"`
	CrearSQL   bool `json:"crearsql"`
	Todo       bool `json:"todo"`
	Prueba     bool `json:"prueba"`
	Hack       bool `json:"hack"`
	Desarrollo bool `json:"desarrollo"`
	Consola    bool `json:"consola"`
	Funcion    bool `json:"funcion"`
}

// Privilegio
type Privilegio struct {
	Id          string   `json:"id"`
	Controlador string   `json:"controlador"`
	Metodo      string   `json:"metodo"`
	Accion      string   `json:"accion"`
	Parametros  []string `json:"parametros"`
}

// Perfil
type Perfil struct {
	Id          string       `json:"id"`
	Descripcion string       `json:"descripcion"`
	Privilegios []Privilegio `json:"privilegios"`
}

type Rol struct {
	Id           string `json:"id"`
	Descripcion  string `json:"descripcion"`
	MetodoSeguro `json:"metodoseguro"`
}

// Usuarios del Sistema
type Usuario struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
	Correo string `json:"correo,omitempty"`
	Clave  string `json:"clave,omitempty"`
	Perfil `json:"perfil,omitempty"`
	Rol    `json:"rol,omitempty"`
}

// La firma permite identificar una maquina y persona autorizada por el sistema
type FirmaDigital struct {
	Id           int
	Usuario      Usuario
	DireccionMac string
	DireccionIP  string
	Tiempo       time.Time
}

func (f *FirmaDigital) Registrar() bool {

	return true
}
