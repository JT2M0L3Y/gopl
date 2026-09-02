package semantic

import (
	"fmt"
	"gopl/internal/ast"
	"gopl/internal/token"
)

var baseTypes = map[string]struct{}{
	"int": {}, "double": {}, "char": {}, "string": {}, "bool": {}, "void": {},
}

// SemanticChecker performs type checking for the GoPL AST.
type SemanticChecker struct {
	ast.Visitor
	SymbolTable *SymbolTable
	CurrType    ast.DataType
	StructDefs  map[string]ast.StructDef
	FunDefs     map[string]ast.FuncDef
}

// NewSemanticChecker creates a new semantic checker.
func NewSemanticChecker() *SemanticChecker {
	return &SemanticChecker{
		SymbolTable: NewSymbolTable(),
		StructDefs:  map[string]ast.StructDef{},
		FunDefs:     map[string]ast.FuncDef{},
	}
}

// Check validates a whole program.
func (sc *SemanticChecker) Check(prog *ast.Program) error {
	for i := range prog.StructDefs {
		name := prog.StructDefs[i].StructName.Lexeme
		if _, exists := sc.StructDefs[name]; exists {
			return sc.errorf("multiple definitions of '%s'", name)
		}
		sc.StructDefs[name] = prog.StructDefs[i]
	}

	foundMain := false
	for i := range prog.FunDefs {
		name := prog.FunDefs[i].FunName.Lexeme
		if _, exists := sc.FunDefs[name]; exists {
			return sc.errorf("multiple definitions of '%s'", name)
		}
		if name == "main" {
			foundMain = true
			if !sameType(prog.FunDefs[i].ReturnType, ast.DataType{TypeNames: []string{"void"}}) {
				return sc.errorf("main function must have void type")
			}
			if len(prog.FunDefs[i].Params) != 0 {
				return sc.errorf("main function cannot have parameters")
			}
		}
		sc.FunDefs[name] = prog.FunDefs[i]
	}
	if !foundMain {
		return sc.errorf("program missing main function")
	}

	for i := range prog.StructDefs {
		if err := prog.StructDefs[i].Accept(sc); err != nil {
			return err
		}
	}
	for i := range prog.FunDefs {
		if err := prog.FunDefs[i].Accept(sc); err != nil {
			return err
		}
	}
	return nil
}

func (sc *SemanticChecker) errorf(msg string, args ...interface{}) error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return fmt.Errorf("semantic error: %s", msg)
}

func (sc *SemanticChecker) isBaseType(dt ast.DataType) bool {
	if len(dt.TypeNames) == 0 {
		return false
	}
	_, ok := baseTypes[dt.TypeNames[0]]
	return ok
}

func (sc *SemanticChecker) isStructType(dt ast.DataType) bool {
	if len(dt.TypeNames) == 0 {
		return false
	}
	_, ok := sc.StructDefs[dt.TypeNames[0]]
	return ok
}

func (sc *SemanticChecker) typeValid(dt ast.DataType) bool {
	if len(dt.TypeNames) == 0 {
		return false
	}
	if sc.isBaseType(dt) || sc.isStructType(dt) {
		return true
	}
	return false
}

func (sc *SemanticChecker) resolveTypeForToken(tok token.Token) string {
	if tok.Kind == token.IntVal {
		return "int"
	}
	if tok.Kind == token.FloatVal {
		return "double"
	}
	if tok.Kind == token.CharVal {
		return "char"
	}
	if tok.Kind == token.StringVal {
		return "string"
	}
	if tok.Kind == token.BoolVal {
		return "bool"
	}
	if tok.Kind == token.NullVal {
		return "void"
	}
	return tok.Lexeme
}

func sameType(a, b ast.DataType) bool {
	if a.IsArray != b.IsArray {
		return false
	}
	if len(a.TypeNames) != len(b.TypeNames) {
		return false
	}
	for i := range a.TypeNames {
		if a.TypeNames[i] != b.TypeNames[i] {
			return false
		}
	}
	return true
}

func assignCompatible(lhs, rhs ast.DataType) bool {
	if sameType(lhs, rhs) {
		return true
	}
	if lhs.IsArray && rhs.IsArray && len(lhs.TypeNames) == 1 && len(rhs.TypeNames) == 1 && lhs.TypeNames[0] == rhs.TypeNames[0] {
		return true
	}
	return false
}

// Visitor methods

func (sc *SemanticChecker) VisitProgram(p *ast.Program) error {
	return sc.Check(p)
}

func (sc *SemanticChecker) VisitFuncDef(f *ast.FuncDef) error {
	sc.SymbolTable.PushEnvironment()
	defer sc.SymbolTable.PopEnvironment()

	if !sc.typeValid(f.ReturnType) {
		return sc.errorf("function return type cannot be '%s'", f.ReturnType.TypeNames[0])
	}
	if f.FunName.Lexeme == "main" && len(f.Params) != 0 {
		return sc.errorf("main function cannot have parameters")
	}

	sc.SymbolTable.Add("return", f.ReturnType)
	for i := range f.Params {
		name := f.Params[i].VarName.Lexeme
		if sc.SymbolTable.NameExistsInCurrEnv(name) {
			return sc.errorf("multiple definitions of '%s'", name)
		}
		if !sc.typeValid(f.Params[i].DataType) {
			return sc.errorf("function parameter cannot have type '%s'", f.Params[i].DataType.TypeNames[0])
		}
		sc.SymbolTable.Add(name, f.Params[i].DataType)
	}

	for i := range f.Stmts {
		if err := f.Stmts[i].Accept(sc); err != nil {
			return err
		}
	}
	return nil
}

func (sc *SemanticChecker) VisitStructDef(s *ast.StructDef) error {
	sc.SymbolTable.PushEnvironment()
	defer sc.SymbolTable.PopEnvironment()

	for i := range s.Fields {
		field := s.Fields[i]
		if sc.SymbolTable.NameExistsInCurrEnv(field.VarName.Lexeme) {
			return sc.errorf("multiple definitions of '%s'", field.VarName.Lexeme)
		}
		if !sc.typeValid(field.DataType) {
			return sc.errorf("struct field cannot have type '%s'", field.DataType.TypeNames[0])
		}
		sc.SymbolTable.Add(field.VarName.Lexeme, field.DataType)
	}
	return nil
}

func (sc *SemanticChecker) VisitReturnStmt(r *ast.ReturnStmt) error {
	if err := r.Expr.Accept(sc); err != nil {
		return err
	}
	retType, ok := sc.SymbolTable.Get("return")
	if !ok {
		return sc.errorf("return type not found")
	}
	if !assignCompatible(retType, sc.CurrType) {
		return sc.errorf("type mismatch in return statement")
	}
	return nil
}

func (sc *SemanticChecker) VisitWhileStmt(w *ast.WhileStmt) error {
	if err := w.Condition.Accept(sc); err != nil {
		return err
	}
	if !sameType(sc.CurrType, ast.DataType{TypeNames: []string{"bool"}}) {
		return sc.errorf("while condition must be of type bool")
	}

	sc.SymbolTable.PushEnvironment()
	defer sc.SymbolTable.PopEnvironment()
	for i := range w.Stmts {
		if err := w.Stmts[i].Accept(sc); err != nil {
			return err
		}
	}
	return nil
}

func (sc *SemanticChecker) VisitForStmt(f *ast.ForStmt) error {
	sc.SymbolTable.PushEnvironment()
	defer sc.SymbolTable.PopEnvironment()
	if !sameType(f.VarDecl.VarDef.DataType, ast.DataType{TypeNames: []string{"int"}}) {
		return sc.errorf("for iterator must be integer type")
	}
	if err := f.VarDecl.Accept(sc); err != nil {
		return err
	}
	sc.SymbolTable.Add(f.VarDecl.VarDef.VarName.Lexeme, f.VarDecl.VarDef.DataType)
	if err := f.Condition.Accept(sc); err != nil {
		return err
	}
	if !sameType(sc.CurrType, ast.DataType{TypeNames: []string{"bool"}}) {
		return sc.errorf("condition must be of bool type")
	}
	if err := f.Assign.Accept(sc); err != nil {
		return err
	}
	for i := range f.Stmts {
		if err := f.Stmts[i].Accept(sc); err != nil {
			return err
		}
	}
	return nil
}

func (sc *SemanticChecker) VisitIfStmt(i *ast.IfStmt) error {
	if err := i.IfPart.Condition.Accept(sc); err != nil {
		return err
	}
	if !sameType(sc.CurrType, ast.DataType{TypeNames: []string{"bool"}}) {
		return sc.errorf("if condition must be of type bool")
	}

	sc.SymbolTable.PushEnvironment()
	for j := range i.IfPart.Stmts {
		if err := i.IfPart.Stmts[j].Accept(sc); err != nil {
			return err
		}
	}
	sc.SymbolTable.PopEnvironment()

	for j := range i.ElseIfs {
		if err := i.ElseIfs[j].Condition.Accept(sc); err != nil {
			return err
		}
		if !sameType(sc.CurrType, ast.DataType{TypeNames: []string{"bool"}}) {
			return sc.errorf("elseif condition must be of type bool")
		}
		sc.SymbolTable.PushEnvironment()
		for k := range i.ElseIfs[j].Stmts {
			if err := i.ElseIfs[j].Stmts[k].Accept(sc); err != nil {
				return err
			}
		}
		sc.SymbolTable.PopEnvironment()
	}

	sc.SymbolTable.PushEnvironment()
	for j := range i.ElseStmts {
		if err := i.ElseStmts[j].Accept(sc); err != nil {
			return err
		}
	}
	sc.SymbolTable.PopEnvironment()
	return nil
}

func (sc *SemanticChecker) VisitVarDeclStmt(v *ast.VarDeclStmt) error {
	if sc.SymbolTable.NameExistsInCurrEnv(v.VarDef.VarName.Lexeme) {
		return sc.errorf("multiple definitions of '%s'", v.VarDef.VarName.Lexeme)
	}
	if !sc.typeValid(v.VarDef.DataType) {
		return sc.errorf("undefined type in var decl")
	}
	if err := v.Expr.Accept(sc); err != nil {
		return err
	}
	if !assignCompatible(v.VarDef.DataType, sc.CurrType) {
		return sc.errorf("type mismatch in var decl")
	}
	sc.SymbolTable.Add(v.VarDef.VarName.Lexeme, v.VarDef.DataType)
	return nil
}

func (sc *SemanticChecker) VisitAssignStmt(a *ast.AssignStmt) error {
	if len(a.LValue) == 0 {
		return nil
	}
	name := a.LValue[0].VarName.Lexeme
	lhsType, ok := sc.SymbolTable.Get(name)
	if !ok {
		return sc.errorf("variable '%s' not defined", name)
	}
	if err := a.Expr.Accept(sc); err != nil {
		return err
	}
	if !assignCompatible(lhsType, sc.CurrType) {
		return sc.errorf("types do not match in assignment")
	}
	return nil
}

func (sc *SemanticChecker) VisitCallExpr(c *ast.CallExpr) error {
	fname := c.FunName.Lexeme
	if fname == "print" {
		if len(c.Args) != 1 {
			return sc.errorf("print expects 1 argument")
		}
		if err := c.Args[0].Accept(sc); err != nil {
			return err
		}
		sc.CurrType = ast.DataType{TypeNames: []string{"void"}}
		return nil
	}
	if fname == "input" {
		if len(c.Args) != 0 {
			return sc.errorf("input does not expect arguments")
		}
		sc.CurrType = ast.DataType{TypeNames: []string{"string"}}
		return nil
	}
	if fname == "length" {
		if len(c.Args) != 1 {
			return sc.errorf("length expects 1 argument")
		}
		if err := c.Args[0].Accept(sc); err != nil {
			return err
		}
		sc.CurrType = ast.DataType{TypeNames: []string{"int"}}
		return nil
	}
	if fname == "to_string" || fname == "to_int" || fname == "to_double" {
		if len(c.Args) != 1 {
			return sc.errorf("%s expects 1 argument", fname)
		}
		if err := c.Args[0].Accept(sc); err != nil {
			return err
		}
		if fname == "to_string" {
			sc.CurrType = ast.DataType{TypeNames: []string{"string"}}
			return nil
		}
		if fname == "to_int" {
			sc.CurrType = ast.DataType{TypeNames: []string{"int"}}
			return nil
		}
		sc.CurrType = ast.DataType{TypeNames: []string{"double"}}
		return nil
	}

	fun, ok := sc.FunDefs[fname]
	if !ok {
		return sc.errorf("function '%s' was not defined", fname)
	}
	if len(c.Args) != len(fun.Params) {
		return sc.errorf("function '%s' expects %d arguments", fname, len(fun.Params))
	}
	for i := range c.Args {
		if err := c.Args[i].Accept(sc); err != nil {
			return err
		}
		if !assignCompatible(fun.Params[i].DataType, sc.CurrType) {
			return sc.errorf("function '%s' expects argument %d to be of type %s", fname, i+1, fun.Params[i].DataType.TypeNames[0])
		}
	}
	sc.CurrType = fun.ReturnType
	return nil
}

func (sc *SemanticChecker) VisitExpr(e *ast.Expr) error {
	if e == nil {
		sc.CurrType = ast.DataType{TypeNames: []string{"void"}}
		return nil
	}

	if e.Op == nil {
		if e.First == nil {
			sc.CurrType = ast.DataType{TypeNames: []string{"void"}}
			return nil
		}
		if err := e.First.Accept(sc); err != nil {
			return err
		}
		if e.Negated {
			if !sameType(sc.CurrType, ast.DataType{TypeNames: []string{"bool"}}) {
				return sc.errorf("negating non-bool type '%s'", sc.CurrType.TypeNames[0])
			}
		}
		return nil
	}

	if err := e.First.Accept(sc); err != nil {
		return err
	}
	leftType := sc.CurrType
	if err := e.Rest.Accept(sc); err != nil {
		return err
	}
	rightType := sc.CurrType

	if e.Negated {
		if !sameType(leftType, ast.DataType{TypeNames: []string{"bool"}}) {
			return sc.errorf("negating non-bool type '%s'", leftType.TypeNames[0])
		}
	}

	switch e.Op.Kind {
	case token.And, token.Or:
		if !sameType(leftType, ast.DataType{TypeNames: []string{"bool"}}) || !sameType(rightType, ast.DataType{TypeNames: []string{"bool"}}) {
			return sc.errorf("using and/or ops on non-bool types")
		}
		sc.CurrType = ast.DataType{TypeNames: []string{"bool"}}
	case token.Plus, token.Minus, token.Times, token.Divide:
		if (!sameType(leftType, ast.DataType{TypeNames: []string{"int"}}) || !sameType(rightType, ast.DataType{TypeNames: []string{"int"}})) &&
			(!sameType(leftType, ast.DataType{TypeNames: []string{"double"}}) || !sameType(rightType, ast.DataType{TypeNames: []string{"double"}})) {
			return sc.errorf("using arithmetic ops on non-int/double types")
		}
		sc.CurrType = leftType
	case token.Equal, token.NotEqual, token.Less, token.Greater, token.LessEq, token.GreaterEq:
		if !sameType(leftType, rightType) {
			return sc.errorf("using relational/equality ops on mismatched types")
		}
		sc.CurrType = ast.DataType{TypeNames: []string{"bool"}}
	default:
		sc.CurrType = leftType
	}
	return nil
}

func (sc *SemanticChecker) VisitSimpleTerm(t *ast.SimpleTerm) error {
	if t == nil || t.RValue == nil {
		sc.CurrType = ast.DataType{TypeNames: []string{"void"}}
		return nil
	}
	return t.RValue.Accept(sc)
}

func (sc *SemanticChecker) VisitComplexTerm(t *ast.ComplexTerm) error {
	if err := t.Expr.Accept(sc); err != nil {
		return err
	}
	return nil
}

func (sc *SemanticChecker) VisitSimpleRValue(v *ast.SimpleRValue) error {
	if v == nil {
		sc.CurrType = ast.DataType{TypeNames: []string{"void"}}
		return nil
	}
	sc.CurrType = ast.DataType{TypeNames: []string{sc.resolveTypeForToken(v.Value)}}
	return nil
}

func (sc *SemanticChecker) VisitNewRValue(v *ast.NewRValue) error {
	if v.ArrayExpr != nil {
		if err := v.ArrayExpr.Accept(sc); err != nil {
			return err
		}
		if !sameType(sc.CurrType, ast.DataType{TypeNames: []string{"int"}}) {
			return sc.errorf("new array size must be int")
		}
		sc.CurrType = ast.DataType{IsArray: true, TypeNames: []string{v.Type.Lexeme}}
		return nil
	}
	sc.CurrType = ast.DataType{TypeNames: []string{v.Type.Lexeme}}
	return nil
}

func (sc *SemanticChecker) VisitVarRValue(v *ast.VarRValue) error {
	if len(v.Path) == 0 {
		sc.CurrType = ast.DataType{TypeNames: []string{"void"}}
		return nil
	}
	name := v.Path[0].VarName.Lexeme
	if !sc.SymbolTable.NameExists(name) {
		return sc.errorf("variable '%s' not defined", name)
	}
	info, _ := sc.SymbolTable.Get(name)
	sc.CurrType = info
	return nil
}
