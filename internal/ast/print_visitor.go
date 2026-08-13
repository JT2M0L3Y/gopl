package ast

import (
	"fmt"
	"io"
	"strings"

	"gopl/internal/token"
)

// PrintVisitor implements the Visitor interface for pretty-printing the AST
type PrintVisitor struct {
	out       io.Writer
	indent    int
	indentAmt int
}

// NewPrintVisitor creates a new PrintVisitor that writes to the given writer
func NewPrintVisitor(w io.Writer) *PrintVisitor {
	return &PrintVisitor{
		out:       w,
		indent:    0,
		indentAmt: 2,
	}
}

// Helper methods

func (pv *PrintVisitor) incIndent() {
	pv.indent += pv.indentAmt
}

func (pv *PrintVisitor) decIndent() {
	pv.indent -= pv.indentAmt
}

func (pv *PrintVisitor) printIndent() {
	fmt.Fprint(pv.out, strings.Repeat(" ", pv.indent))
}

func (pv *PrintVisitor) printf(format string, args ...interface{}) {
	fmt.Fprintf(pv.out, format, args...)
}

// Visitor interface implementation

func (pv *PrintVisitor) VisitProgram(p *Program) error {
	for _, structDef := range p.StructDefs {
		if err := structDef.Accept(pv); err != nil {
			return err
		}
	}
	for _, funDef := range p.FunDefs {
		if err := funDef.Accept(pv); err != nil {
			return err
		}
	}
	return nil
}

func (pv *PrintVisitor) VisitFuncDef(f *FuncDef) error {
	pv.printf("\n")

	// Print return type
	if f.ReturnType.IsArray {
		pv.printf("array %s ", f.ReturnType.TypeNames[0])
	} else if f.ReturnType.IsDict {
		pv.printf("dict %s %s ", f.ReturnType.TypeNames[0], f.ReturnType.TypeNames[1])
	} else {
		pv.printf("%s ", f.ReturnType.TypeNames[0])
	}

	// Print function name and parameters
	pv.printf("%s(", f.FunName.Lexeme)

	for i, param := range f.Params {
		if param.DataType.IsArray {
			pv.printf("array %s ", param.DataType.TypeNames[0])
		} else if param.DataType.IsDict {
			pv.printf("dict %s %s ", param.DataType.TypeNames[0], param.DataType.TypeNames[1])
		} else {
			pv.printf("%s ", param.DataType.TypeNames[0])
		}

		if i != len(f.Params)-1 {
			pv.printf("%s, ", param.VarName.Lexeme)
		} else {
			pv.printf("%s", param.VarName.Lexeme)
		}
	}
	pv.printf(") {\n")

	// Print function body
	pv.incIndent()
	for _, stmt := range f.Stmts {
		pv.printIndent()
		if err := stmt.Accept(pv); err != nil {
			return err
		}
		pv.printf("\n")
	}
	pv.decIndent()
	pv.printIndent()
	pv.printf("}\n")

	return nil
}

func (pv *PrintVisitor) VisitStructDef(s *StructDef) error {
	pv.printf("\nstruct %s {\n", s.StructName.Lexeme)

	pv.incIndent()
	for i, field := range s.Fields {
		pv.printIndent()

		if field.DataType.IsArray {
			pv.printf("array %s ", field.DataType.TypeNames[0])
		} else if field.DataType.IsDict {
			pv.printf("dict %s %s ", field.DataType.TypeNames[0], field.DataType.TypeNames[1])
		} else {
			pv.printf("%s ", field.DataType.TypeNames[0])
		}

		if i == len(s.Fields)-1 {
			pv.printf("%s\n", field.VarName.Lexeme)
		} else {
			pv.printf("%s,\n", field.VarName.Lexeme)
		}
	}
	pv.decIndent()
	pv.printf("}\n")

	return nil
}

func (pv *PrintVisitor) VisitReturnStmt(s *ReturnStmt) error {
	pv.printf("return ")
	return s.Expr.Accept(pv)
}

func (pv *PrintVisitor) VisitWhileStmt(s *WhileStmt) error {
	pv.printf("while ")
	if err := s.Condition.Accept(pv); err != nil {
		return err
	}
	pv.printf(" {\n")

	pv.incIndent()
	for _, stmt := range s.Stmts {
		pv.printIndent()
		if err := stmt.Accept(pv); err != nil {
			return err
		}
		pv.printf("\n")
	}
	pv.decIndent()
	pv.printIndent()
	pv.printf("}")

	return nil
}

func (pv *PrintVisitor) VisitForStmt(s *ForStmt) error {
	pv.printf("for (")
	if err := s.VarDecl.Accept(pv); err != nil {
		return err
	}
	pv.printf("; ")
	if err := s.Condition.Accept(pv); err != nil {
		return err
	}
	pv.printf("; ")
	if err := s.Assign.Accept(pv); err != nil {
		return err
	}
	pv.printf(") {\n")

	pv.incIndent()
	for _, stmt := range s.Stmts {
		pv.printIndent()
		if err := stmt.Accept(pv); err != nil {
			return err
		}
		pv.printf("\n")
	}
	pv.decIndent()
	pv.printIndent()
	pv.printf("}")

	return nil
}

func (pv *PrintVisitor) VisitIfStmt(s *IfStmt) error {
	pv.printf("if (")
	if err := s.IfPart.Condition.Accept(pv); err != nil {
		return err
	}
	pv.printf(") {\n")

	pv.incIndent()
	for _, stmt := range s.IfPart.Stmts {
		pv.printIndent()
		if err := stmt.Accept(pv); err != nil {
			return err
		}
		pv.printf("\n")
	}
	pv.decIndent()

	// Handle elseif statements
	if len(s.ElseIfs) > 0 {
		for _, bi := range s.ElseIfs {
			pv.printIndent()
			pv.printf("}\n")
			pv.printIndent()
			pv.printf("elseif (")
			if err := bi.Condition.Accept(pv); err != nil {
				return err
			}
			pv.printf(") {\n")
			pv.incIndent()

			for _, stmt := range bi.Stmts {
				pv.printIndent()
				if err := stmt.Accept(pv); err != nil {
					return err
				}
				pv.printf("\n")
			}
			pv.decIndent()
		}
	}

	// Handle else statement
	if len(s.ElseStmts) > 0 {
		pv.printIndent()
		pv.printf("}\n")
		pv.printIndent()
		pv.printf("else {\n")
		pv.incIndent()

		for _, stmt := range s.ElseStmts {
			pv.printIndent()
			if err := stmt.Accept(pv); err != nil {
				return err
			}
			pv.printf("\n")
		}
		pv.decIndent()
	}

	pv.printIndent()
	pv.printf("}")

	return nil
}

func (pv *PrintVisitor) VisitVarDeclStmt(s *VarDeclStmt) error {
	if s.VarDef.DataType.IsArray {
		pv.printf("array %s ", s.VarDef.DataType.TypeNames[0])
	} else if s.VarDef.DataType.IsDict {
		pv.printf("dict %s %s ", s.VarDef.DataType.TypeNames[0], s.VarDef.DataType.TypeNames[1])
	} else {
		pv.printf("%s ", s.VarDef.DataType.TypeNames[0])
	}

	pv.printf("%s = ", s.VarDef.VarName.Lexeme)
	return s.Expr.Accept(pv)
}

func (pv *PrintVisitor) VisitAssignStmt(s *AssignStmt) error {
	if len(s.LValue) == 1 {
		pv.printf("%s", s.LValue[0].VarName.Lexeme)
		if s.LValue[0].ArrayExpr != nil {
			pv.printf("[")
			if err := s.LValue[0].ArrayExpr.Accept(pv); err != nil {
				return err
			}
			pv.printf("]")
		} else if s.LValue[0].DictExpr != nil {
			pv.printf("[")
			if err := s.LValue[0].DictExpr.Accept(pv); err != nil {
				return err
			}
			pv.printf("]")
		}
	} else {
		for i, lval := range s.LValue {
			pv.printf("%s", lval.VarName.Lexeme)
			if lval.ArrayExpr != nil {
				pv.printf("[")
				if err := lval.ArrayExpr.Accept(pv); err != nil {
					return err
				}
				pv.printf("]")
			} else if lval.DictExpr != nil {
				pv.printf("[")
				if err := lval.DictExpr.Accept(pv); err != nil {
					return err
				}
				pv.printf("]")
			}
			if i != len(s.LValue)-1 {
				pv.printf(".")
			}
		}
	}

	pv.printf(" = ")
	return s.Expr.Accept(pv)
}

func (pv *PrintVisitor) VisitCallExpr(e *CallExpr) error {
	pv.printf("%s(", e.FunName.Lexeme)

	if len(e.Args) == 1 {
		if err := e.Args[0].Accept(pv); err != nil {
			return err
		}
	} else if len(e.Args) > 1 {
		for i, arg := range e.Args {
			if err := arg.Accept(pv); err != nil {
				return err
			}
			if i != len(e.Args)-1 {
				pv.printf(", ")
			}
		}
	}

	pv.printf(")")
	return nil
}

func (pv *PrintVisitor) VisitExpr(e *Expr) error {
	if e.Op != nil {
		if e.Negated {
			pv.printf("not (")
			if err := e.First.Accept(pv); err != nil {
				return err
			}
			pv.printf(" %s ", e.Op.Lexeme)
			if err := e.Rest.Accept(pv); err != nil {
				return err
			}
			pv.printf(")")
		} else {
			if err := e.First.Accept(pv); err != nil {
				return err
			}
			pv.printf(" %s ", e.Op.Lexeme)
			if err := e.Rest.Accept(pv); err != nil {
				return err
			}
		}
	} else {
		if e.Negated {
			pv.printf("not (")
			if err := e.First.Accept(pv); err != nil {
				return err
			}
			pv.printf(")")
		} else {
			if err := e.First.Accept(pv); err != nil {
				return err
			}
		}
	}

	return nil
}

func (pv *PrintVisitor) VisitSimpleTerm(t *SimpleTerm) error {
	return t.RValue.Accept(pv)
}

func (pv *PrintVisitor) VisitComplexTerm(t *ComplexTerm) error {
	pv.printf("(")
	if err := t.Expr.Accept(pv); err != nil {
		return err
	}
	pv.printf(")")
	return nil
}

func (pv *PrintVisitor) VisitSimpleRValue(v *SimpleRValue) error {
	switch v.Value.Kind {
	case token.StringVal:
		pv.printf("\"%s\"", v.Value.Lexeme)
	case token.CharVal:
		pv.printf("'%s'", v.Value.Lexeme)
	default:
		pv.printf("%s", v.Value.Lexeme)
	}
	return nil
}

func (pv *PrintVisitor) VisitNewRValue(v *NewRValue) error {
	pv.printf("new %s", v.Type.Lexeme)

	if v.ArrayExpr != nil {
		pv.printf("[")
		if err := v.ArrayExpr.Accept(pv); err != nil {
			return err
		}
		pv.printf("]")
	} else if v.DictExpr != nil {
		pv.printf("{")
		if err := v.DictExpr.Accept(pv); err != nil {
			return err
		}
		pv.printf("}")
	}

	return nil
}

func (pv *PrintVisitor) VisitVarRValue(v *VarRValue) error {
	if len(v.Path) == 1 {
		pv.printf("%s", v.Path[0].VarName.Lexeme)
		if v.Path[0].ArrayExpr != nil {
			pv.printf("[")
			if err := v.Path[0].ArrayExpr.Accept(pv); err != nil {
				return err
			}
			pv.printf("]")
		} else if v.Path[0].DictExpr != nil {
			pv.printf("[")
			if err := v.Path[0].DictExpr.Accept(pv); err != nil {
				return err
			}
			pv.printf("]")
		}
	} else {
		for i, varRef := range v.Path {
			if i == 0 {
				pv.printf("%s", varRef.VarName.Lexeme)
			} else {
				pv.printf(".%s", varRef.VarName.Lexeme)
			}

			if varRef.ArrayExpr != nil {
				pv.printf("[")
				if err := varRef.ArrayExpr.Accept(pv); err != nil {
					return err
				}
				pv.printf("]")
			} else if varRef.DictExpr != nil {
				pv.printf("[")
				if err := varRef.DictExpr.Accept(pv); err != nil {
					return err
				}
				pv.printf("]")
			}
		}
	}

	return nil
}
