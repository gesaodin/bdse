//ubicación geográfica de los modelos
package localizacion

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gesaodin/bdse/sys"
)

//Estado es la forma de organización política,
//dotada de poder soberano e independiente, que integra la
//población de un territorio
type Estado struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

//Ciudad s un asentamiento de población con atribuciones y
//funciones político-administrativas, económicas y religiosas,
//a diferencia de los núcleos rurales
type Ciudad struct {
	ID       int    `json:"id"`
	IDEstado int    `json:"ide,omitempty"`
	Nombre   string `json:"nombre"`
}

//Municipio  está compuesto por un territorio claramente
//definido por un término municipal de límites fijados
type Municipio struct {
	ID       int    `json:"id"`
	IDEstado int    `json:"ide,omitempty"`
	Nombre   string `json:"nombre"`
}

//Parroquia es la denominación de algunas entidades subnacionales
//en diferentes países
type Parroquia struct {
	ID          int    `json:"id"`
	IDMunicipio int    `json:"idm,omitempty"`
	Nombre      string `json:"nombre"`
}

//Respuesta Generales
type Respuesta struct {
	Cantidad int64  `json:"cant"` // Cantidad de elementos
	Msj      string `json:"msj"`  // Mensaje almacenado
}

//Consultar una lista de todos los estados
func (e *Estado) Consultar() (jSon []byte, err error) {
	var donde string
	if e.ID > 0 {
		donde = ` WHERE id_estado = ` + strconv.Itoa(e.ID) + `;`
	}
	s := `SELECT id_estado,estado FROM estados ` + donde
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}

	var lst []interface{}
	for row.Next() {
		var estado Estado
		var id int
		var nombre string
		e := row.Scan(&id, &nombre)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		estado.ID = id
		estado.Nombre = nombre
		lst = append(lst, estado)
	}

	jSon, _ = json.Marshal(lst)
	return

}

//Consultar una lista de todos los Ciudad
func (c *Ciudad) Consultar() (jSon []byte, err error) {
	var donde string
	if c.IDEstado > 0 {
		donde = ` WHERE id_estado = ` + strconv.Itoa(c.IDEstado) + `;`
	}
	s := `SELECT id_ciudad,ciudad FROM ciudades ` + donde
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}

	var lst []interface{}
	for row.Next() {
		var ciudad Ciudad
		var id int
		var nombre string
		e := row.Scan(&id, &nombre)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		ciudad.ID = id
		ciudad.Nombre = nombre
		lst = append(lst, ciudad)
	}

	jSon, _ = json.Marshal(lst)
	return
}

//Consultar una lista de todos los Municipio
func (m *Municipio) Consultar() (jSon []byte, err error) {
	var donde string
	if m.IDEstado > 0 {
		donde = ` WHERE id_estado = ` + strconv.Itoa(m.IDEstado) + `;`
	}
	s := `SELECT id_municipio,municipio FROM municipios ` + donde
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}

	var lst []interface{}
	for row.Next() {
		var municipio Municipio
		var id int
		var nombre string
		e := row.Scan(&id, &nombre)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		municipio.ID = id
		municipio.Nombre = nombre
		lst = append(lst, municipio)
	}

	jSon, _ = json.Marshal(lst)
	return
}

//Consultar una lista de todos los Municipio
func (p *Parroquia) Consultar() (jSon []byte, err error) {
	var donde string

	if p.IDMunicipio > 0 {
		donde = ` WHERE id_municipio = ` + strconv.Itoa(p.IDMunicipio) + `;`
	}
	s := `SELECT id_parroquia,parroquia FROM parroquias ` + donde
	row, err := sys.PostgreSQL.Query(s)
	if err != nil {
		return
	}

	var lst []interface{}
	for row.Next() {
		var parroquia Parroquia
		var id int
		var nombre string
		e := row.Scan(&id, &nombre)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		parroquia.ID = id
		parroquia.Nombre = nombre
		lst = append(lst, parroquia)
	}

	jSon, _ = json.Marshal(lst)
	return
}
