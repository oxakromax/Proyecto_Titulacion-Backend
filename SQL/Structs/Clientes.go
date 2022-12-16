package Structs

import "sync"

type Clientes struct {
	ID             int    `json:"ID"`
	Nombre         string `json:"Nombre"`
	Apellido       string `json:"Apellido"`
	Email          string `json:"Email"`
	OrganizacionID int    `json:"OrganizacionID"`
}

type ClientesArray struct {
	Clientes []Clientes `json:"Clientes"`
	mu       sync.Mutex
}

func (c *ClientesArray) Set(value []Clientes) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Clientes = value
}

func (c *ClientesArray) Get() []Clientes {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Clientes
}

func (c *ClientesArray) Add(value Clientes) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Clientes = append(c.Clientes, value)
}

func (c *ClientesArray) Delete(value Clientes) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, v := range c.Clientes {
		if v == value {
			c.Clientes = append(c.Clientes[:i], c.Clientes[i+1:]...)
		}
	}
}
