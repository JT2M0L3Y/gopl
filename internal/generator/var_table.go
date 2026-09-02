package generator

import "slices"

type VarTable struct {
	environments []map[string]int
}

func NewVarTable() *VarTable {
	return &VarTable{environments: []map[string]int{{}}}
}

func (t *VarTable) PushEnvironment() {
	t.environments = append(t.environments, map[string]int{})
}

func (t *VarTable) PopEnvironment() {
	if len(t.environments) > 1 {
		t.environments = t.environments[:len(t.environments)-1]
	}
}

func (t *VarTable) Add(name string) int {
	i := t.Size()
	t.environments[len(t.environments)-1][name] = i
	return i
}

func (t *VarTable) Size() int {
	n := 0
	for _, env := range t.environments {
		n += len(env)
	}
	return n
}

func (t *VarTable) Get(name string) (int, bool) {
	for _, v := range slices.Backward(t.environments) {
		if n, ok := v[name]; ok {
			return n, true
		}
	}
	return 0, false
}
