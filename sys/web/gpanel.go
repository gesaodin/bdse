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
	DescripcionDelPanel string = "Bus de Servicio Empresarial"
	VersionDelPanel     string = "V.0.0.1"
	AutorDelPanel       string = "Carlos E. Peña A."
)

type Pagina struct {
	Urlcss string
	Urljs  string
}

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

type WebData struct {
	Title string
}

func (G *GPanel) Login(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	//fmt.Println(r.RemoteAddr)
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
	case "login":
		G.TituloDePagina = "login"
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
			G.TituloDePagina = "login"
			G.IrA(w)
		} //Err Session
	} //Fin switch
}

func (G *GPanel) Validar(w http.ResponseWriter, r *http.Request) {
	var usuario seguridad.Usuario

	session, e := seguridad.Stores.Get(r, "session-bdse")
	if e != nil {
		G.TextoError = e.Error()
		G.Error(w)
		//http.Redirect(w, r, "login", http.StatusFound)
		return
	}

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
		G.TituloDePagina = "login"
		G.IrA(w)
	}

}

func (G *GPanel) IrA(w http.ResponseWriter) {
	// fmt.Println("Entrando en funcion ", G.TituloDePagina)
	var t *template.Template
	var err error
	var base string = "public_web/AdminLTE/"

	if G.TituloDePagina != "login" {

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
		valor := strings.Split(strings.Split(strings.Trim(cadena[1], " "), ".")[0], " ")
		fecha = valor[2] + "-" + valor[1] + "-" + valor[0]
		tipoArchivo(fecha, files[i].Filename, usuario, codigo)

	}
}

//En caso de acceder a una url sin acceso
func (G *GPanel) Error(w http.ResponseWriter) {
	terr, _ := template.ParseFiles("public_web/AdminLTE/err.html")
	terr.Execute(w, G)

}

//Salir del Panel o finalizar sesión
func (G *GPanel) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := seguridad.Stores.Get(r, "session-bdse")
	session.Values["acceso"] = false
	sessions.Save(r, w)
	G.TituloDePagina = "login"
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
	archivo.Ch = Mensajeria.Usuario[usuario].ch
	archivo.Fecha = f

	// t := strings.Split(s, ".")
	// switch strings.ToLower(t[0]) {
	switch strings.ToLower(codigo) {
	case "ma":
		go archivo.LeerMaticlo(Mensajeria.Usuario["gpanel"].ch)
		return
	case "mo":
		go archivo.LeerMorpheus(Mensajeria.Usuario["gpanel"].ch)
		return
	case "p1":
		go archivo.LeerPos(Mensajeria.Usuario["gpanel"].ch, 2)
		return
	case "p2":
		go archivo.LeerPos(Mensajeria.Usuario["gpanel"].ch, 3)
		return
	case "p3":
		go archivo.LeerPos(Mensajeria.Usuario["gpanel"].ch, 4)
		return
	case "sp":
		go archivo.LeerSport(Mensajeria.Usuario["gpanel"].ch)
		return
	case "il":
		go archivo.LeerIlbanquero(Mensajeria.Usuario["gpanel"].ch)
		return
	case "cy":
		go archivo.LeerCyberParley(Mensajeria.Usuario["gpanel"].ch)
		return
	case "a1":
		go archivo.LeerPos(Mensajeria.Usuario["gpanel"].ch, 9)
		return
	case "a2":
		go archivo.LeerPos(Mensajeria.Usuario["gpanel"].ch, 10)
		return
	case "a3":
		go archivo.LeerPos(Mensajeria.Usuario["gpanel"].ch, 11)
		return
	case "t1":
		go archivo.LeerPos(Mensajeria.Usuario["gpanel"].ch, 12)
		return
	default:

		return
	}

}

func (G *GPanel) Reporte(w http.ResponseWriter, r *http.Request) {

}
