package ast

import (
	"testing"

	"gopl/internal/token"
)

type recordingVisitor struct {
	visited string
}

func (v *recordingVisitor) VisitProgram(*Program) error       { v.visited = "program"; return nil }
func (v *recordingVisitor) VisitFuncDef(*FuncDef) error       { v.visited = "func"; return nil }
func (v *recordingVisitor) VisitStructDef(*StructDef) error   { v.visited = "struct"; return nil }
func (v *recordingVisitor) VisitReturnStmt(*ReturnStmt) error { v.visited = "return"; return nil }
func (v *recordingVisitor) VisitWhileStmt(*WhileStmt) error   { v.visited = "while"; return nil }
func (v *recordingVisitor) VisitForStmt(*ForStmt) error       { v.visited = "for"; return nil }
func (v *recordingVisitor) VisitIfStmt(*IfStmt) error         { v.visited = "if"; return nil }
func (v *recordingVisitor) VisitVarDeclStmt(*VarDeclStmt) error {
	v.visited = "var-decl"
	return nil
}
func (v *recordingVisitor) VisitAssignStmt(*AssignStmt) error { v.visited = "assign"; return nil }
func (v *recordingVisitor) VisitCallExpr(*CallExpr) error     { v.visited = "call"; return nil }
func (v *recordingVisitor) VisitExpr(*Expr) error             { v.visited = "expr"; return nil }
func (v *recordingVisitor) VisitSimpleTerm(*SimpleTerm) error { v.visited = "simple-term"; return nil }
func (v *recordingVisitor) VisitComplexTerm(*ComplexTerm) error {
	v.visited = "complex-term"
	return nil
}
func (v *recordingVisitor) VisitSimpleRValue(*SimpleRValue) error {
	v.visited = "simple-rvalue"
	return nil
}
func (v *recordingVisitor) VisitNewRValue(*NewRValue) error {
	v.visited = "new-rvalue"
	return nil
}
func (v *recordingVisitor) VisitVarRValue(*VarRValue) error {
	v.visited = "var-rvalue"
	return nil
}

func TestAcceptDispatchesToMatchingVisitorMethod(t *testing.T) {
	cases := []struct {
		name string
		node ASTNode
		want string
	}{
		{"program", &Program{}, "program"},
		{"function", &FuncDef{}, "func"},
		{"struct", &StructDef{}, "struct"},
		{"return", &ReturnStmt{}, "return"},
		{"while", &WhileStmt{}, "while"},
		{"for", &ForStmt{}, "for"},
		{"if", &IfStmt{}, "if"},
		{"variable declaration", &VarDeclStmt{}, "var-decl"},
		{"assignment", &AssignStmt{}, "assign"},
		{"call", &CallExpr{}, "call"},
		{"expression", &Expr{}, "expr"},
		{"simple term", &SimpleTerm{}, "simple-term"},
		{"complex term", &ComplexTerm{}, "complex-term"},
		{"simple rvalue", &SimpleRValue{}, "simple-rvalue"},
		{"new rvalue", &NewRValue{}, "new-rvalue"},
		{"variable rvalue", &VarRValue{}, "var-rvalue"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			visitor := &recordingVisitor{}
			if err := test.node.Accept(visitor); err != nil {
				t.Fatal(err)
			}
			if visitor.visited != test.want {
				t.Fatalf("visited %q, want %q", visitor.visited, test.want)
			}
		})
	}
}

func TestFirstTokenPropagation(t *testing.T) {
	first := token.Token{Kind: token.IntVal, Lexeme: "42", Line: 4, Col: 9}
	second := token.Token{Kind: token.Ident, Lexeme: "value", Line: 4, Col: 12}
	value := &SimpleRValue{Value: first}
	term := &SimpleTerm{RValue: value}
	expr := &Expr{First: term}
	complexTerm := &ComplexTerm{Expr: *expr}
	call := &CallExpr{FunName: second}
	variable := &VarRValue{Path: []VarRef{{VarName: second}}}

	cases := []struct {
		name string
		node interface{ FirstToken() token.Token }
		want token.Token
	}{
		{"variable definition", VarDef{VarName: second}, second},
		{"simple value", value, first},
		{"simple term", term, first},
		{"expression", expr, first},
		{"complex term", complexTerm, first},
		{"call", call, second},
		{"variable value", variable, second},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if got := test.node.FirstToken(); got != test.want {
				t.Fatalf("token = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestFirstTokenForEmptyNodesIsZeroToken(t *testing.T) {
	if got := (&Expr{}).FirstToken(); got != (token.Token{}) {
		t.Fatalf("empty expression token = %#v", got)
	}
	if got := (&SimpleTerm{}).FirstToken(); got != (token.Token{}) {
		t.Fatalf("empty term token = %#v", got)
	}
	if got := (&VarRValue{}).FirstToken(); got != (token.Token{}) {
		t.Fatalf("empty variable value token = %#v", got)
	}
}
