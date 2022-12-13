package Structs

import (
	"sync"
	"time"
)

// create table public.incidentes_procesos
//(
//    id         serial
//        primary key,
//    proceso_id integer           not null
//        references public.procesos,
//    incidente  text,
//    tipo       integer default 1 not null,
//    estado     integer default 1 not null
//);
//
//alter table public.incidentes_procesos
//    owner to postgres;

//create table public.incidentes_detalle
//(
//    id           serial
//        primary key,
//    incidente_id integer not null
//        references public.incidentes_procesos,
//    detalle      text,
//    fecha_inicio timestamp,
//    fecha_fin    timestamp
//);
//
//alter table public.incidentes_detalle
//    owner to postgres;

type Incidentes_Procesos struct {
	ID        int    `json:"ID"`
	ProcesoID int    `json:"ProcesoID"`
	Incidente string `json:"Incidente"`
	Tipo      int    `json:"Tipo"`
	Estado    int    `json:"Estado"`
}

type Incidentes_ProcesosArray struct {
	Incidentes_Procesos []Incidentes_Procesos `json:"Incidentes_Procesos"`
	mu                  sync.Mutex
}

func (i *Incidentes_ProcesosArray) Set(value []Incidentes_Procesos) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Incidentes_Procesos = value
}

func (i *Incidentes_ProcesosArray) Get() []Incidentes_Procesos {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.Incidentes_Procesos
}

func (i *Incidentes_ProcesosArray) Add(value Incidentes_Procesos) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Incidentes_Procesos = append(i.Incidentes_Procesos, value)
}

func (i *Incidentes_ProcesosArray) Delete(value Incidentes_Procesos) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for index, v := range i.Incidentes_Procesos {
		if v == value {
			i.Incidentes_Procesos = append(i.Incidentes_Procesos[:index], i.Incidentes_Procesos[index+1:]...)
		}
	}
}

type Incidentes_Detalle struct {
	ID           int       `json:"ID"`
	IncidenteID  int       `json:"IncidenteID"`
	Detalle      string    `json:"Detalle"`
	Fecha_Inicio time.Time `json:"Fecha_Inicio"`
	Fecha_Fin    time.Time `json:"Fecha_Fin"`
}

type Incidentes_DetalleArray struct {
	Incidentes_Detalle []Incidentes_Detalle `json:"Incidentes_Detalle"`
	mu                 sync.Mutex
}

func (i *Incidentes_DetalleArray) Set(value []Incidentes_Detalle) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Incidentes_Detalle = value
}

func (i *Incidentes_DetalleArray) Get() []Incidentes_Detalle {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.Incidentes_Detalle
}

func (i *Incidentes_DetalleArray) Add(value Incidentes_Detalle) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Incidentes_Detalle = append(i.Incidentes_Detalle, value)
}

func (i *Incidentes_DetalleArray) Delete(value Incidentes_Detalle) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for index, v := range i.Incidentes_Detalle {
		if v == value {
			i.Incidentes_Detalle = append(i.Incidentes_Detalle[:index], i.Incidentes_Detalle[index+1:]...)
		}
	}
}
