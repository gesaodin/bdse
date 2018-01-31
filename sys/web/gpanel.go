//sistemas de paginas y archivos para el www
package web

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gesaodin/bdse/sys"
	"github.com/gesaodin/bdse/sys/seguridad"
	"github.com/gesaodin/bdse/util"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	//DescripcionDelPanel Descripcion
	DescripcionDelPanel string = "Bus de Servicio Empresarial"
	//VersionDelPanel Version
	VersionDelPanel string = "V.0.0.1"
	//AutorDelPanel Autor
	AutorDelPanel string = "Carlos E. Peña A."
	//_Login Usuario
	_Login = "login"
)

//Pagina Direcciones de Cabaeceras
type Pagina struct {
	Urlcss string
	Urljs  string
}

//GPanel Reglas de descripcion general del panel
type GPanel struct {
	Descripcion    string
	Version        string
	Autor          string
	Fecha          time.Time
	Nivel          int
	Pagina         string
	TituloDePagina string
	TextoError     string
	Usuario        seguridad.Usuario
	Config         Pagina
}

//Data Titulo general
type Data struct {
	Title string
}

//Login Inicio de sesion del panel
func (G *GPanel) Login(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	session, e := seguridad.Stores.Get(r, "session-bdse")
	if e != nil {

		G.TextoError = e.Error()
		G.Error(w)
		//http.Redirect(w, r, "login", http.StatusFound)
		return
	}

	if session.Values["acceso"] != nil && session.Values["acceso"].(bool) {
		G.Usuario.Nombre = session.Values["usuario"].(string)
		G.Usuario.Rol = session.Values["rol"].(string)
	}

	G.TituloDePagina = v["id"]
	switch G.TituloDePagina {
	case "validar":
		G.Validar(w, r)
	case "logout":
		G.Logout(w, r)
	case "subirl": //Subir Archivos de Loteria
		G.SubirArchivoLoteria(w, r)
	case "subirp": //Subir Archivos de Parley
		G.SubirArchivoLoteria(w, r)
	case _Login:
		G.TituloDePagina = _Login
		if session.Values["acceso"] != nil && session.Values["acceso"].(bool) {
			G.TituloDePagina = "principal"
		}

		G.IrA(w)
	default:
		if session.Values["acceso"] != nil {
			if session.Values["acceso"].(bool) {
				// fmt.Println("Ruta nueva conectado...")

				G.IrA(w)
			} else {
				G.TextoError = "Acceso denegado cookies caducada"
				G.Error(w)

			} //Session == true
		} else {
			G.TextoError = "Acceso denegado"
			G.TituloDePagina = _Login
			G.IrA(w)
		} //Err Session
	} //Fin switch
}

//Validar Verificacion del usuario
func (G *GPanel) Validar(w http.ResponseWriter, r *http.Request) {
	var usuario seguridad.Usuario

	session, e := seguridad.Stores.Get(r, "session-bdse")
	if e != nil {
		G.TextoError = e.Error()
		G.Error(w)
		//http.Redirect(w, r, _Login, http.StatusFound)
		return
	}
	//fmt.Println(r.FormValue("usuario")  + r.FormValue("clave") )
	if r.FormValue("usuario") != "" {

		b := usuario.Consultar(r.FormValue("usuario"), r.FormValue("clave"))

		if b {
			session.Values["acceso"] = true
			session.Values["usuario"] = r.FormValue("usuario")
			session.Values["rol"] = usuario.Rol
			sessions.Save(r, w)
			G.TituloDePagina = "principal"
			G.Descripcion = DescripcionDelPanel
			G.Version = VersionDelPanel
			G.Autor = AutorDelPanel
			G.Usuario = usuario
			G.IrA(w)

		} else {
			session.Values["rol"] = ""
			session.Values["acceso"] = false
			G.TextoError = "El usuario no se encuentra registrado"
			G.Error(w)
		}

	} else {
		G.TituloDePagina = _Login
		G.IrA(w)
	}

}

func (G *GPanel) IrA(w http.ResponseWriter) {
	// fmt.Println("Entrando en funcion ", G.TituloDePagina)
	var t *template.Template
	var err error
	var base string = "public_web/adminlte/"

	if G.TituloDePagina != _Login {

		plantilla := base + "p" + G.TituloDePagina + ".ghtm"

		G.Config.Urlcss = ".css"
		base += "rol/" + strings.ToLower(G.Usuario.Rol) + "/"
		cabecera := base + "inc/cabecera.html"
		menu := base + "inc/menu.html"
		cuerpo := base + "pag/" + G.TituloDePagina + ".html"
		pie := base + "inc/pie.html"
		t, err = template.ParseFiles(plantilla, cabecera, menu, cuerpo, pie)
		if err != nil {
			G.TextoError = "La painga no se encuentra disponible"
			G.Error(w)
			return
		}
		t.ExecuteTemplate(w, "plantilla", &G)
		t.ExecuteTemplate(w, "cabecera", &G)
		t.ExecuteTemplate(w, "menu", &G)
		t.ExecuteTemplate(w, "contenido", &G)
		t.ExecuteTemplate(w, "pie", &G)
	} else {
		t, err = template.ParseFiles(base + G.TituloDePagina + ".html")
	}

	if err != nil {
		G.Error(w)
	} else {
		t.Execute(w, &G)
	}
}

//Subir archivos al sistema
func (G *GPanel) SubirArchivoLoteria(w http.ResponseWriter, r *http.Request) {

	session, e := seguridad.Stores.Get(r, "session-bdse")

	if e != nil {
		fmt.Println("Error Cookies: ", e)
	}
	er := r.ParseMultipartForm(32 << 20)
	if er != nil {
		fmt.Println(er)
		return
	}
	m := r.MultipartForm
	files := m.File["archivo"]
	fecha := r.FormValue("fecha")
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		// out, err := os.Create("./public/temp/loteria/" + fecha + files[i].Filename)
		out, err := os.Create("./public/temp/loteria/" + files[i].Filename)
		defer out.Close()
		if err != nil {
			fmt.Fprintf(w, "No se pudo escribir el archivo por favor verifique los privilegios.")
			return
		}
		_, err = io.Copy(out, file) // file not files[i] !
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintf(w, "Archivo "+files[i].Filename+" enviado..."+"\n")
		usuario := session.Values["usuario"].(string)
		cadena := strings.Split(files[i].Filename, "-")
		codigo := strings.Trim(cadena[0], " ")
		valor := strings.Split(cadena[3], ".")

		fecha = valor[0] + "-" + cadena[2] + "-" + cadena[1]
		tipoArchivo(fecha, files[i].Filename, usuario, codigo)

	}
}

//En caso de acceder a una url sin acceso
func (G *GPanel) Error(w http.ResponseWriter) {
	terr, _ := template.ParseFiles("public_web/adminlte/err.html")
	terr.Execute(w, G)

}

//Logout del Panel o finalizar sesión
func (G *GPanel) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := seguridad.Stores.Get(r, "session-bdse")
	session.Values["acceso"] = false
	sessions.Save(r, w)
	G.TituloDePagina = _Login
	G.IrA(w)
}

//Identificar el archivo que se está cargado
func tipoArchivo(f string, s string, usuario string, codigo string) {

	var archivo = util.Archivo{}
	archivo.Ruta = "public/temp/loteria/"
	// archivo.NombreDelArchivo = f + s
	archivo.NombreDelArchivo = s
	sys.PostgreSQL.SetMaxOpenConns(10)
	archivo.PostgreSQL = sys.PostgreSQL
	archivo.Canal = Mensajeria.Usuario["gpanel"].ch
	archivo.Fecha = f
	tipo := strings.ToLower(codigo[2:3])
	tipopos := strings.ToLower(codigo[1:3])

	switch strings.ToLower(codigo[:2]) {
	case "ma": //Maticlot
		go archivo.LeerMaticlo(Mensajeria.Usuario[usuario].ch, tipo)
		break
	case "mo": //Morpheus
		go archivo.LeerMorpheus(Mensajeria.Usuario[usuario].ch, tipo)
		break
	case "p1": //POS
		go archivo.LeerPos(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "p2":
		go archivo.LeerPos(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "p3":
		go archivo.LeerPos(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "a1": //Aliens
		go archivo.LeerAliens(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "a2":
		go archivo.LeerAliens(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "a3":
		go archivo.LeerAliens(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "t1": //Turco
		go archivo.LeerTurco(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "t2":
		go archivo.LeerTurco(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "t3":
		go archivo.LeerTurco(Mensajeria.Usuario[usuario].ch, tipopos)
		break
	case "il":
		go archivo.LeerIlbanquero(Mensajeria.Usuario[usuario].ch, tipo)
		break
	case "mx":
		go archivo.LeerMatrix(Mensajeria.Usuario[usuario].ch, tipo)
		break
	case "cy":
		go archivo.LeerCyberParley(Mensajeria.Usuario[usuario].ch, tipo)
		break
	case "sp":
		go archivo.LeerSport(Mensajeria.Usuario[usuario].ch, tipo)
		break
	case "mp": //Morpheus
		go archivo.LeerMatchPoint(Mensajeria.Usuario[usuario].ch, tipo)
		break
	default:

		break
	}

}

func (G *GPanel) Reporte(w http.ResponseWriter, r *http.Request) {

}
