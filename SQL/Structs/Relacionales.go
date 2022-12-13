package Structs

import "sync"

// create table public.procesos_clientes
//(
//    id         serial
//        primary key,
//    proceso_id integer not null
//        references public.procesos,
//    cliente_id integer not null
//        references public.clientes
//);
//
//alter table public.procesos_clientes
//    owner to postgres;

type Procesos_Clientes struct {
	ID        int `json:"ID"`
	ProcesoID int `json:"ProcesoID"`
	ClienteID int `json:"ClienteID"`
}

type Procesos_ClientesArray struct {
	Procesos_Clientes []Procesos_Clientes `json:"Procesos_Clientes"`
	mu                sync.Mutex
}

func (p *Procesos_ClientesArray) Set(value []Procesos_Clientes) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Procesos_Clientes = value
}

func (p *Procesos_ClientesArray) Get() []Procesos_Clientes {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Procesos_Clientes
}

func (p *Procesos_ClientesArray) Add(value Procesos_Clientes) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Procesos_Clientes = append(p.Procesos_Clientes, value)
}

func (p *Procesos_ClientesArray) Delete(value Procesos_Clientes) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, v := range p.Procesos_Clientes {
		if v == value {
			p.Procesos_Clientes = append(p.Procesos_Clientes[:i], p.Procesos_Clientes[i+1:]...)
		}
	}
}

// create table public.procesos_usuarios
//(
//    id         serial
//        primary key,
//    proceso_id integer not null
//        references public.procesos,
//    usuario_id integer not null
//        references public.usuarios
//);
//
//alter table public.procesos_usuarios
//    owner to postgres;

type Procesos_Usuarios struct {
	ID        int `json:"ID"`
	ProcesoID int `json:"ProcesoID"`
	UsuarioID int `json:"UsuarioID"`
}

type Procesos_UsuariosArray struct {
	Procesos_Usuarios []Procesos_Usuarios `json:"Procesos_Usuarios"`
	mu                sync.Mutex
}

func (p *Procesos_UsuariosArray) Set(value []Procesos_Usuarios) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Procesos_Usuarios = value

}

func (p *Procesos_UsuariosArray) Get() []Procesos_Usuarios {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Procesos_Usuarios
}

func (p *Procesos_UsuariosArray) Add(value Procesos_Usuarios) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Procesos_Usuarios = append(p.Procesos_Usuarios, value)
}

func (p *Procesos_UsuariosArray) Delete(value Procesos_Usuarios) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, v := range p.Procesos_Usuarios {
		if v == value {
			p.Procesos_Usuarios = append(p.Procesos_Usuarios[:i], p.Procesos_Usuarios[i+1:]...)
		}
	}
}
