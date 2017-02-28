package agencia

type Agencia struct {
	Oid         int         `json:"oid,omitempty"`
	Nombre      string      `json:"nombre,omitempty"`
	Telefono    string      `json:"telefono,omitempty"`
	Fecha       string      `json:"fecha,omitempty"`
	Observacion string			`json:"observacion,omitempty"`
	Responsable interface{} `json:"responsable,omitempty"`
	Taquilla    interface{} `json:"taquilla,omitempty"`
	Saldo       float32     `json:"saldo,omitempty"`
	Direccion   interface{} `json:"direccion,omitempty"`
}

type Responsable struct {
	Oid       int         `json:"oid,omitempty"`
	Cedula    string      `json:"cedula,omitempty"`
	Nombre    string      `json:"nombre,omitempty"`
	Telefono  string      `json:"telefono,omitempty"`
	Direccion interface{} `json:"direccion,omitempty"`
}

type Taquilla struct {
	Oid    int    `json:"oid,omitempty"`
	Nombre string `json:"nombre,omitempty"`
	Fecha  string `json:"fecha,omitempty"`
}

type Direccion struct {
	Ciudad    int    `json:"ciudad,omitempty"`
	Estado    int    `json:"estado,omitempty"`
	Municipio int    `json:"municipio,omitempty"`
	Parroquia int    `json:"parroquia,omitempty"`
	Sector    string `json:"sector,omitempty"`
	Casa      string `json:"casa,omitempty"`
}

type Sistema struct {
	Oid       int    `json:"oid,omitempty"`
	IdSistema int    `json:"idsistema,omitempty"`
	Fecha     string `json:"fecha,omitempty"`
}

//Registrar una Agencia
func (a *Agencia) Registrar (agencia Agencia) (jSon []byte, err error) {

	return
}

//Consultar una Agencia
func (a *Agencia) Consultar (agencia Agencia) (jSon []byte, err error) {

	return
}
