package ast

import "gopl/internal/token"

//----------------------------------------------------------------------
// Abstract AST Node Interfaces
//----------------------------------------------------------------------

// ASTNode is the base interface for all AST nodes
type ASTNode interface {
	Accept(Visitor) error
}

// Stmt is the interface for all statement nodes
type Stmt interface {
	ASTNode
}

// ExprTerm is the interface for expression terms
type ExprTerm interface {
	ASTNode
	FirstToken() token.Token
}

// RValue is the interface for right-hand values
type RValue interface {
	ASTNode
	FirstToken() token.Token
}

//----------------------------------------------------------------------
// Visitor interface
//----------------------------------------------------------------------

type Visitor interface {
	// top-level
	VisitProgram(*Program) error
	VisitFuncDef(*FuncDef) error
	VisitStructDef(*StructDef) error
	// statements
	VisitReturnStmt(*ReturnStmt) error
	VisitWhileStmt(*WhileStmt) error
	VisitForStmt(*ForStmt) error
	VisitIfStmt(*IfStmt) error
	VisitVarDeclStmt(*VarDeclStmt) error
	VisitAssignStmt(*AssignStmt) error
	VisitCallExpr(*CallExpr) error
	VisitExpr(*Expr) error
	VisitSimpleTerm(*SimpleTerm) error
	VisitComplexTerm(*ComplexTerm) error
	VisitSimpleRValue(*SimpleRValue) error
	VisitNewRValue(*NewRValue) error
	VisitVarRValue(*VarRValue) error
}

//----------------------------------------------------------------------
// Program-related types
//----------------------------------------------------------------------

type Program struct {
	ASTNode
	StructDefs []StructDef
	FunDefs    []FuncDef
}

func (p *Program) Accept(v Visitor) error {
	return v.VisitProgram(p)
}

type DataType struct {
	IsArray   bool
	IsDict    bool
	TypeNames []string
}

type VarDef struct {
	DataType DataType
	VarName  token.Token
}

func (v VarDef) FirstToken() token.Token {
	return v.VarName
}

type StructDef struct {
	ASTNode
	StructName token.Token
	Fields     []VarDef
}

func (s *StructDef) Accept(v Visitor) error {
	return v.VisitStructDef(s)
}

type FuncDef struct {
	ASTNode
	ReturnType DataType
	FunName    token.Token
	Params     []VarDef
	Stmts      []Stmt
}

func (f *FuncDef) Accept(v Visitor) error {
	return v.VisitFuncDef(f)
}

//----------------------------------------------------------------------
// Expression-related types
//----------------------------------------------------------------------

type Expr struct {
	ASTNode
	Negated bool
	First   ExprTerm
	Op      *token.Token
	Rest    *Expr
}

func (e *Expr) Accept(v Visitor) error {
	return v.VisitExpr(e)
}

func (e *Expr) FirstToken() token.Token {
	if e.First != nil {
		return e.First.FirstToken()
	}
	return token.Token{}
}

type SimpleTerm struct {
	ExprTerm
	RValue RValue
}

func (s *SimpleTerm) Accept(v Visitor) error {
	return v.VisitSimpleTerm(s)
}

func (s *SimpleTerm) FirstToken() token.Token {
	if s.RValue != nil {
		return s.RValue.FirstToken()
	}
	return token.Token{}
}

type ComplexTerm struct {
	ExprTerm
	Expr Expr
}

func (c *ComplexTerm) Accept(v Visitor) error {
	return v.VisitComplexTerm(c)
}

func (c *ComplexTerm) FirstToken() token.Token {
	return c.Expr.FirstToken()
}

type SimpleRValue struct {
	RValue
	Value token.Token
}

func (s *SimpleRValue) Accept(v Visitor) error {
	return v.VisitSimpleRValue(s)
}

func (s *SimpleRValue) FirstToken() token.Token {
	return s.Value
}

type NewRValue struct {
	RValue
	Type      token.Token
	ArrayExpr *Expr
	DictExpr  *Expr
}

func (n *NewRValue) Accept(v Visitor) error {
	return v.VisitNewRValue(n)
}

func (n *NewRValue) FirstToken() token.Token {
	return n.Type
}

type VarRef struct {
	VarName   token.Token
	ArrayExpr *Expr
	DictExpr  *Expr
}

type VarRValue struct {
	RValue
	Path []VarRef
}

func (v *VarRValue) Accept(visitor Visitor) error {
	return visitor.VisitVarRValue(v)
}

func (v *VarRValue) FirstToken() token.Token {
	if len(v.Path) > 0 {
		return v.Path[0].VarName
	}
	return token.Token{}
}

//----------------------------------------------------------------------
// Statement-related types
//----------------------------------------------------------------------

type ReturnStmt struct {
	Stmt
	Expr Expr
}

func (r *ReturnStmt) Accept(v Visitor) error {
	return v.VisitReturnStmt(r)
}

type WhileStmt struct {
	Stmt
	Condition Expr
	Stmts     []Stmt
}

func (w *WhileStmt) Accept(v Visitor) error {
	return v.VisitWhileStmt(w)
}

type VarDeclStmt struct {
	Stmt
	VarDef VarDef
	Expr   Expr
}

func (v *VarDeclStmt) Accept(visitor Visitor) error {
	return visitor.VisitVarDeclStmt(v)
}

type AssignStmt struct {
	Stmt
	LValue []VarRef
	Expr   Expr
}

func (a *AssignStmt) Accept(v Visitor) error {
	return v.VisitAssignStmt(a)
}

type ForStmt struct {
	Stmt
	VarDecl   VarDeclStmt
	Condition Expr
	Assign    AssignStmt
	Stmts     []Stmt
}

func (f *ForStmt) Accept(v Visitor) error {
	return v.VisitForStmt(f)
}

type BasicIf struct {
	Condition Expr
	Stmts     []Stmt
}

type IfStmt struct {
	Stmt
	IfPart    BasicIf
	ElseIfs   []BasicIf
	ElseStmts []Stmt
}

func (i *IfStmt) Accept(v Visitor) error {
	return v.VisitIfStmt(i)
}

type CallExpr struct {
	Stmt
	RValue
	FunName token.Token
	Args    []Expr
}

func (c *CallExpr) Accept(v Visitor) error {
	return v.VisitCallExpr(c)
}

func (c *CallExpr) FirstToken() token.Token {
	return c.FunName
}
