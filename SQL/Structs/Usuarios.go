package Structs

type AuthUsuario struct {
	Email    string
	Password string
}

type Usuario struct {
	ID       int    `json:"ID"`
	Nombre   string `json:"Nombre"`
	Apellido string `json:"Apellido"`
	Email    string `json:"Email"`
	Roles    []struct {
		ID     int    `json:"ID"`
		Nombre string `json:"Nombre"`
	}
}
