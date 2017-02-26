package sys

import (
	"database/sql"
	"encoding/json"

	mgo "gopkg.in/mgo.v2"

	"github.com/gesaodin/bdse/util"
)

type config struct{}

var (
	Version string = "V.0.0.1"
	// PostgreSQL  bool   = false
	MySQL       bool = false
	MongoDB     bool = false
	SQLServer   bool = false
	Oracle      bool = false
	BaseDeDatos BaseDatos
	MGOSession  *mgo.Session
	PostgreSQL  *sql.DB
	Error       error
)

const (
	ActivarConexionRemota      bool   = true
	DesactivarConexionRemota   bool   = false
	ActivarLogDeRegistro       bool   = true
	DesactivarLogDeRegistro    bool   = false
	ActivarRoles               bool   = true
	DesactivarRoles            bool   = false
	ActivarLimitEnConsultas    bool   = true
	DesactivarLimitEnConsultas bool   = false
	Puerto                     string = "2004"
	PuertoSSL                  string = "2608"
	CodificacionDeArchivo      string = "UTF-8"
	MaximoLimiteDeUsuarios     int    = 100
	MaximoLimiteDeConsultas    int    = 10
)

type BaseDatos struct {
	CadenaDeConexion map[string]CadenaDeConexion
}

type CadenaDeConexion struct {
	Driver    string
	Usuario   string
	Basedatos string
	Clave     string
	Host      string
	Puerto    string
}

//0: PostgreSQL, 1: MySQL, 2: MongoDB
var Conexiones []CadenaDeConexion

func init() {
	var a util.Archivo
	a.NombreDelArchivo = "sys/config.json"
	data, _ := a.LeerTodo()
	e := json.Unmarshal(data, &Conexiones)
	for _, valor := range Conexiones {
		switch valor.Driver {
		case "postgres":
			cad := make(map[string]CadenaDeConexion)
			cad["postgres"] = CadenaDeConexion{
				Driver:    valor.Driver,
				Usuario:   valor.Usuario,
				Basedatos: valor.Basedatos,
				Clave:     valor.Clave,
				Host:      valor.Host,
				Puerto:    valor.Puerto,
			}
			BaseDeDatos.CadenaDeConexion = cad
		case "mysql":
			MySQL = true
		case "mongodb":
			MongoDB = true
		}
	}
	util.Error(e)
}
