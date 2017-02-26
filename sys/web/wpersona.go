package web

//Manejador de Contenido URL
func WMPersona() {

	// Enrutador.HandleFunc("/admin/validar", ValidarTokenNew)
	Enrutador.HandleFunc("/persona", ValidarToken(PersonaGET)) //Obtener
	//Enrutador.HandleFunc("/persona", PersonaPOST).Methods("POST")                  //Agregar
	//Enrutador.HandleFunc("/persona/{id}", PersonaUpdate).Methods("PUT", "OPTIONS") //Actualizar
	//Enrutador.HandleFunc("/persona", PersonaUpdate) //.Methods("OPTIONS") //Actualizar
	//Enrutador.HandleFunc("/personaUpdate", PersonaPUT)                      //Actualizar
}
