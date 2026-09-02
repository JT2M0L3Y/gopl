package semantic

import "gopl/internal/ast"

// SymbolTable stores type information in nested scopes.
type SymbolTable struct {
	environments []map[string]ast.DataType
}

// NewSymbolTable creates an empty symbol table with a global scope.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		environments: []map[string]ast.DataType{{}},
	}
}

// PushEnvironment adds a new scope.
func (st *SymbolTable) PushEnvironment() {
	st.environments = append(st.environments, map[string]ast.DataType{})
}

// PopEnvironment removes the most recent scope.
func (st *SymbolTable) PopEnvironment() {
	if len(st.environments) > 1 {
		st.environments = st.environments[:len(st.environments)-1]
	}
}

// Empty reports whether there are no environments.
func (st *SymbolTable) Empty() bool {
	return len(st.environments) == 0
}

// Add inserts a symbol into the current environment.
func (st *SymbolTable) Add(name string, info ast.DataType) {
	st.environments[len(st.environments)-1][name] = info
}

// NameExists reports whether a symbol exists in any visible scope.
func (st *SymbolTable) NameExists(name string) bool {
	for i := len(st.environments) - 1; i >= 0; i-- {
		if _, ok := st.environments[i][name]; ok {
			return true
		}
	}
	return false
}

// NameExistsInCurrEnv reports whether a symbol exists in the current scope.
func (st *SymbolTable) NameExistsInCurrEnv(name string) bool {
	_, ok := st.environments[len(st.environments)-1][name]
	return ok
}

// Get returns the most recent matching symbol type.
func (st *SymbolTable) Get(name string) (ast.DataType, bool) {
	for i := len(st.environments) - 1; i >= 0; i-- {
		if v, ok := st.environments[i][name]; ok {
			return v, true
		}
	}
	return ast.DataType{}, false
}
