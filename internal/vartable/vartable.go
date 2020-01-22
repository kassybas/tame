package vartable

import (
	"sync"

	"github.com/kassybas/tame/internal/tvar"
)

type VarTable struct {
	sync.RWMutex
	vars map[string]tvar.TVariable
}

func CopyVarTable(vt *VarTable) *VarTable {
	newVt := VarTable{
		vars: make(map[string]tvar.TVariable),
	}
	vt.RLock()
	for k, v := range vt.vars {
		newVt.vars[k] = v
	}
	vt.RUnlock()
	return &newVt
}

func NewVarTable() *VarTable {
	vt := VarTable{}
	vt.vars = make(map[string]tvar.TVariable)
	return &vt
}

func (vt *VarTable) Exists(name string) bool {
	vt.RLock()
	_, exists := vt.vars[name]
	vt.RUnlock()
	return exists
}
