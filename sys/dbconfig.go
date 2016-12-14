package sys

type BaseDeDatosPermisos struct {
	CrearBaseDeDatos     bool
	CrearTablas          bool
	CrearFunciones       bool
	CrearDisparadores    bool
	EliminarBaseDeDatos  bool
	EliminarTablas       bool
	EliminarFunciones    bool
	EliminarDisparadores bool
}
