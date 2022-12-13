package Structs

import (
	"golang.org/x/exp/slices"
	"sync"
)

type ProcesosBDD struct {
	ID               int    `json:"ID"`
	Nombre           string `json:"Nombre"`
	Folderid         int    `json:"Folderid"`
	WarningTolerance int    `json:"WarningTolerance"`
	ErrorTolerance   int    `json:"ErrorTolerance"`
	FatalTolerance   int    `json:"FatalTolerance"`
	OrganizacionId   int    `json:"OrganizacionId"`
}

type ProcessBDDArray struct {
	Processes []ProcesosBDD `json:"Processes"`
	mu        sync.Mutex
}

func (p *ProcessBDDArray) Set(value []ProcesosBDD) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Processes = value
}

func (p *ProcessBDDArray) Get() []ProcesosBDD {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Processes
}

func (p *ProcessBDDArray) Add(value ProcesosBDD) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Processes = append(p.Processes, value)
}

func (p *ProcessBDDArray) Delete(value ProcesosBDD) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, v := range p.Processes {
		if v == value {
			p.Processes = append(p.Processes[:i], p.Processes[i+1:]...)
		}
	}
}

func (p *ProcessBDDArray) FilterUniqueFoldersID() []int {
	p.mu.Lock()
	defer p.mu.Unlock()
	var foldersID []int
	for _, v := range p.Processes {
		if !slices.Contains(foldersID, v.Folderid) {
			foldersID = append(foldersID, v.Folderid)
		}
	}
	return foldersID
}
