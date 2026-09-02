package generator

import (
	"fmt"
	"strconv"
	"strings"

	"gopl/internal/ast"
	"gopl/internal/token"
	"gopl/internal/vm"
)

type Generator struct {
	vm      *vm.VM
	frame   vm.FrameInfo
	vars    *VarTable
	structs map[string]ast.StructDef
}

func New(v *vm.VM) *Generator {
	return &Generator{vm: v, vars: NewVarTable(), structs: map[string]ast.StructDef{}}
}

func (g *Generator) Generate(program *ast.Program) error { return program.Accept(g) }

func (g *Generator) emit(op vm.OpCode, operand interface{}) {
	g.frame.Instructions = append(g.frame.Instructions, vm.NewInstr(op, operand))
}

func (g *Generator) VisitProgram(p *ast.Program) error {
	for i := range p.StructDefs {
		g.structs[p.StructDefs[i].StructName.Lexeme] = p.StructDefs[i]
	}
	for i := range p.FunDefs {
		if err := p.FunDefs[i].Accept(g); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) VisitStructDef(*ast.StructDef) error { return nil }

func (g *Generator) VisitFuncDef(f *ast.FuncDef) error {
	g.frame = vm.FrameInfo{Name: f.FunName.Lexeme, ArgCount: len(f.Params)}
	g.vars = NewVarTable()
	for _, param := range f.Params {
		index := g.vars.Add(param.VarName.Lexeme)
		g.emit(vm.STORE, index)
	}
	for _, stmt := range f.Stmts {
		if err := stmt.Accept(g); err != nil {
			return err
		}
	}
	if len(f.Stmts) == 0 || f.ReturnType.TypeNames[0] == "void" {
		g.emit(vm.PUSH, nil)
		g.emit(vm.RET, nil)
	}
	g.vm.Add(&g.frame)
	return nil
}

func (g *Generator) VisitReturnStmt(s *ast.ReturnStmt) error {
	if err := s.Expr.Accept(g); err != nil {
		return err
	}
	g.emit(vm.RET, nil)
	return nil
}

func (g *Generator) VisitVarDeclStmt(s *ast.VarDeclStmt) error {
	index := g.vars.Add(s.VarDef.VarName.Lexeme)
	if err := s.Expr.Accept(g); err != nil {
		return err
	}
	g.emit(vm.STORE, index)
	return nil
}

func (g *Generator) VisitAssignStmt(s *ast.AssignStmt) error {
	if len(s.LValue) == 0 {
		return fmt.Errorf("assignment has no target")
	}
	first, ok := g.vars.Get(s.LValue[0].VarName.Lexeme)
	if !ok {
		return fmt.Errorf("unknown variable '%s'", s.LValue[0].VarName.Lexeme)
	}
	if len(s.LValue) == 1 && s.LValue[0].ArrayExpr == nil {
		if err := s.Expr.Accept(g); err != nil {
			return err
		}
		g.emit(vm.STORE, first)
		return nil
	}
	g.emit(vm.LOAD, first)
	for i, ref := range s.LValue[:len(s.LValue)-1] {
		if ref.ArrayExpr != nil {
			if err := ref.ArrayExpr.Accept(g); err != nil {
				return err
			}
			if i < len(s.LValue)-1 {
				g.emit(vm.GETI, nil)
			}
		} else if i > 0 {
			g.emit(vm.GETF, ref.VarName.Lexeme)
		}
	}
	if err := s.Expr.Accept(g); err != nil {
		return err
	}
	last := s.LValue[len(s.LValue)-1]
	if last.ArrayExpr != nil {
		if len(s.LValue) > 1 {
			g.emit(vm.GETF, last.VarName.Lexeme)
		}
		g.emit(vm.SETI, nil)
	} else {
		g.emit(vm.SETF, last.VarName.Lexeme)
	}
	return nil
}

func (g *Generator) VisitWhileStmt(s *ast.WhileStmt) error {
	start := len(g.frame.Instructions)
	if err := s.Condition.Accept(g); err != nil {
		return err
	}
	g.emit(vm.JMPF, -1)
	branch := len(g.frame.Instructions) - 1
	g.vars.PushEnvironment()
	for _, stmt := range s.Stmts {
		if err := stmt.Accept(g); err != nil {
			return err
		}
	}
	g.vars.PopEnvironment()
	g.emit(vm.JMP, start)
	g.frame.Instructions[branch].Operand = len(g.frame.Instructions)
	return nil
}

func (g *Generator) VisitForStmt(s *ast.ForStmt) error {
	g.vars.PushEnvironment()
	if err := s.VarDecl.Accept(g); err != nil {
		return err
	}
	start := len(g.frame.Instructions)
	if err := s.Condition.Accept(g); err != nil {
		return err
	}
	g.emit(vm.JMPF, -1)
	branch := len(g.frame.Instructions) - 1
	g.vars.PushEnvironment()
	for _, stmt := range s.Stmts {
		if err := stmt.Accept(g); err != nil {
			return err
		}
	}
	g.vars.PopEnvironment()
	if err := s.Assign.Accept(g); err != nil {
		return err
	}
	g.emit(vm.JMP, start)
	g.frame.Instructions[branch].Operand = len(g.frame.Instructions)
	g.vars.PopEnvironment()
	return nil
}

func (g *Generator) VisitIfStmt(s *ast.IfStmt) error {
	endJumps := []int{}
	if err := s.IfPart.Condition.Accept(g); err != nil {
		return err
	}
	g.emit(vm.JMPF, -1)
	next := len(g.frame.Instructions) - 1
	for _, stmt := range s.IfPart.Stmts {
		if err := stmt.Accept(g); err != nil {
			return err
		}
	}
	if len(s.ElseIfs) > 0 || len(s.ElseStmts) > 0 {
		g.emit(vm.JMP, -1)
		endJumps = append(endJumps, len(g.frame.Instructions)-1)
	}
	g.frame.Instructions[next].Operand = len(g.frame.Instructions)
	for _, part := range s.ElseIfs {
		if err := part.Condition.Accept(g); err != nil {
			return err
		}
		g.emit(vm.JMPF, -1)
		next = len(g.frame.Instructions) - 1
		for _, stmt := range part.Stmts {
			if err := stmt.Accept(g); err != nil {
				return err
			}
		}
		g.emit(vm.JMP, -1)
		endJumps = append(endJumps, len(g.frame.Instructions)-1)
		g.frame.Instructions[next].Operand = len(g.frame.Instructions)
	}
	for _, stmt := range s.ElseStmts {
		if err := stmt.Accept(g); err != nil {
			return err
		}
	}
	for _, jump := range endJumps {
		g.frame.Instructions[jump].Operand = len(g.frame.Instructions)
	}
	return nil
}

func (g *Generator) VisitCallExpr(e *ast.CallExpr) error {
	for _, arg := range e.Args {
		if err := arg.Accept(g); err != nil {
			return err
		}
	}
	name := e.FunName.Lexeme
	builtins := map[string]vm.OpCode{
		"print":        vm.WRITE,
		"input":        vm.READ,
		"to_string":    vm.TOSTR,
		"to_int":       vm.TOINT,
		"to_double":    vm.TODBL,
		"length":       vm.SLEN,
		"length@array": vm.ALEN,
		"get":          vm.GETC,
		"concat":       vm.CONCAT,
	}
	if op, ok := builtins[name]; ok {
		g.emit(op, nil)
	} else {
		g.emit(vm.CALL, name)
	}
	return nil
}

func (g *Generator) VisitExpr(e *ast.Expr) error {
	if e.First != nil {
		if err := e.First.Accept(g); err != nil {
			return err
		}
	}
	if e.Rest != nil {
		if err := e.Rest.Accept(g); err != nil {
			return err
		}
		if e.Op != nil {
			g.emit(binaryOpcode(e.Op.Kind), nil)
		}
	}
	if e.Negated {
		g.emit(vm.NOT, nil)
	}
	return nil
}

func binaryOpcode(k token.Kind) vm.OpCode {
	switch k {
	case token.Plus:
		return vm.ADD
	case token.Minus:
		return vm.SUB
	case token.Times:
		return vm.MUL
	case token.Divide:
		return vm.DIV
	case token.And:
		return vm.AND
	case token.Or:
		return vm.OR
	case token.Less:
		return vm.CMPLT
	case token.LessEq:
		return vm.CMPLE
	case token.Greater:
		return vm.CMPGT
	case token.GreaterEq:
		return vm.CMPGE
	case token.Equal:
		return vm.CMPEQ
	default:
		return vm.CMPNE
	}
}

func (g *Generator) VisitSimpleTerm(t *ast.SimpleTerm) error { return t.RValue.Accept(g) }

func (g *Generator) VisitComplexTerm(t *ast.ComplexTerm) error { return t.Expr.Accept(g) }

func (g *Generator) VisitSimpleRValue(v *ast.SimpleRValue) error {
	value := v.Value.Lexeme
	switch v.Value.Kind {
	case token.IntVal:
		n, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		g.emit(vm.PUSH, n)
	case token.FloatVal:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		g.emit(vm.PUSH, n)
	case token.BoolVal:
		g.emit(vm.PUSH, value == "true")
	case token.NullVal:
		g.emit(vm.PUSH, nil)
	default:
		g.emit(vm.PUSH, strings.Trim(value, "\"'"))
	}
	return nil
}

func (g *Generator) VisitVarRValue(v *ast.VarRValue) error {
	if len(v.Path) == 0 {
		return fmt.Errorf("empty variable path")
	}
	index, ok := g.vars.Get(v.Path[0].VarName.Lexeme)
	if !ok {
		return fmt.Errorf("unknown variable '%s'", v.Path[0].VarName.Lexeme)
	}
	g.emit(vm.LOAD, index)
	for i, ref := range v.Path {
		if i == 0 && ref.ArrayExpr == nil {
			continue
		}
		if ref.ArrayExpr != nil {
			if err := ref.ArrayExpr.Accept(g); err != nil {
				return err
			}
			g.emit(vm.GETI, nil)
		} else {
			g.emit(vm.GETF, ref.VarName.Lexeme)
		}
	}
	return nil
}

func (g *Generator) VisitNewRValue(v *ast.NewRValue) error {
	if v.ArrayExpr != nil {
		if err := v.ArrayExpr.Accept(g); err != nil {
			return err
		}
		g.emit(vm.PUSH, nil)
		g.emit(vm.ALLOCA, nil)
		return nil
	}
	g.emit(vm.ALLOCS, nil)
	if def, ok := g.structs[v.Type.Lexeme]; ok {
		for _, field := range def.Fields {
			g.emit(vm.DUP, nil)
			g.emit(vm.ADDF, field.VarName.Lexeme)
			g.emit(vm.DUP, nil)
			g.emit(vm.PUSH, nil)
			g.emit(vm.SETF, field.VarName.Lexeme)
		}
	}
	return nil
}
