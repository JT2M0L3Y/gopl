package ast

import "gopl/internal/token"

// Program is the root of a GoPL source file.
type Program struct {
	Structs []StructDecl
	Funcs   []FuncDecl
}

// StructDecl declares a structure type.
type StructDecl struct {
	Name   token.Token
	Fields []VarDecl
}

// FuncDecl declares a function.
type FuncDecl struct {
	Name       token.Token
	ReturnType TypeRef
	Params     []VarDecl
	Body       []Stmt
}

// VarDecl declares a typed variable or field.
type VarDecl struct {
	Name token.Token
	Type TypeRef
}

// TypeRef describes a type name and optional collection shape.
type TypeRef struct {
	Name  token.Token
	Array bool
	Dict  bool
}

// Stmt is the marker interface for statements.
type Stmt interface {
	stmtNode()
}

// Expr is the marker interface for expressions.
type Expr interface {
	exprNode()
}
