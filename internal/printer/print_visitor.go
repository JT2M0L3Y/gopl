package printer

import (
	"fmt"
	"io"
	"strings"

	"gopl/internal/ast"
	"gopl/internal/token"
)

// PrintVisitor implements the Visitor interface for pretty-printing the AST
type PrintVisitor struct {
	ast.Visitor
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

func (pv *PrintVisitor) VisitProgram(p *ast.Program) error {
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

func (pv *PrintVisitor) VisitFuncDef(f *ast.FuncDef) error {
	pv.printf("\n")

	// Print return type
	if f.ReturnType.IsArray {
		pv.printf("array %s ", f.ReturnType.TypeNames[0])
	} else {
		pv.printf("%s ", f.ReturnType.TypeNames[0])
	}

	// Print function name and parameters
	pv.printf("%s(", f.FunName.Lexeme)

	for i, param := range f.Params {
		if param.DataType.IsArray {
			pv.printf("array %s ", param.DataType.TypeNames[0])
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

func (pv *PrintVisitor) VisitStructDef(s *ast.StructDef) error {
	pv.printf("\nstruct %s {\n", s.StructName.Lexeme)

	pv.incIndent()
	for i, field := range s.Fields {
		pv.printIndent()

		if field.DataType.IsArray {
			pv.printf("array %s ", field.DataType.TypeNames[0])
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

func (pv *PrintVisitor) VisitReturnStmt(r *ast.ReturnStmt) error {
	pv.printf("return ")
	return r.Expr.Accept(pv)
}

func (pv *PrintVisitor) VisitWhileStmt(w *ast.WhileStmt) error {
	pv.printf("while ")
	if err := w.Condition.Accept(pv); err != nil {
		return err
	}
	pv.printf(" {\n")

	pv.incIndent()
	for _, stmt := range w.Stmts {
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

func (pv *PrintVisitor) VisitForStmt(f *ast.ForStmt) error {
	pv.printf("for (")
	if err := f.VarDecl.Accept(pv); err != nil {
		return err
	}
	pv.printf("; ")
	if err := f.Condition.Accept(pv); err != nil {
		return err
	}
	pv.printf("; ")
	if err := f.Assign.Accept(pv); err != nil {
		return err
	}
	pv.printf(") {\n")

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
	pv.printf("}")

	return nil
}

func (pv *PrintVisitor) VisitIfStmt(i *ast.IfStmt) error {
	pv.printf("if (")
	if err := i.IfPart.Condition.Accept(pv); err != nil {
		return err
	}
	pv.printf(") {\n")

	pv.incIndent()
	for _, stmt := range i.IfPart.Stmts {
		pv.printIndent()
		if err := stmt.Accept(pv); err != nil {
			return err
		}
		pv.printf("\n")
	}
	pv.decIndent()

	// Handle elseif statements
	if len(i.ElseIfs) > 0 {
		for _, bi := range i.ElseIfs {
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
	if len(i.ElseStmts) > 0 {
		pv.printIndent()
		pv.printf("}\n")
		pv.printIndent()
		pv.printf("else {\n")
		pv.incIndent()

		for _, stmt := range i.ElseStmts {
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

func (pv *PrintVisitor) VisitVarDeclStmt(v *ast.VarDeclStmt) error {
	if v.VarDef.DataType.IsArray {
		pv.printf("array %s ", v.VarDef.DataType.TypeNames[0])
	} else {
		pv.printf("%s ", v.VarDef.DataType.TypeNames[0])
	}

	pv.printf("%s = ", v.VarDef.VarName.Lexeme)
	return v.Expr.Accept(pv)
}

func (pv *PrintVisitor) VisitAssignStmt(a *ast.AssignStmt) error {
	if len(a.LValue) == 1 {
		pv.printf("%s", a.LValue[0].VarName.Lexeme)
		if a.LValue[0].ArrayExpr != nil {
			pv.printf("[")
			if err := a.LValue[0].ArrayExpr.Accept(pv); err != nil {
				return err
			}
			pv.printf("]")
		}
	} else {
		for i, lval := range a.LValue {
			pv.printf("%s", lval.VarName.Lexeme)
			if lval.ArrayExpr != nil {
				pv.printf("[")
				if err := lval.ArrayExpr.Accept(pv); err != nil {
					return err
				}
				pv.printf("]")
			}
			if i != len(a.LValue)-1 {
				pv.printf(".")
			}
		}
	}

	pv.printf(" = ")
	return a.Expr.Accept(pv)
}

func (pv *PrintVisitor) VisitCallExpr(c *ast.CallExpr) error {
	pv.printf("%s(", c.FunName.Lexeme)

	if len(c.Args) == 1 {
		if err := c.Args[0].Accept(pv); err != nil {
			return err
		}
	} else if len(c.Args) > 1 {
		for i, arg := range c.Args {
			if err := arg.Accept(pv); err != nil {
				return err
			}
			if i != len(c.Args)-1 {
				pv.printf(", ")
			}
		}
	}

	pv.printf(")")
	return nil
}

func (pv *PrintVisitor) VisitExpr(e *ast.Expr) error {
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

func (pv *PrintVisitor) VisitSimpleTerm(s *ast.SimpleTerm) error {
	return s.RValue.Accept(pv)
}

func (pv *PrintVisitor) VisitComplexTerm(c *ast.ComplexTerm) error {
	pv.printf("(")
	if err := c.Expr.Accept(pv); err != nil {
		return err
	}
	pv.printf(")")
	return nil
}

func (pv *PrintVisitor) VisitSimpleRValue(s *ast.SimpleRValue) error {
	switch s.Value.Kind {
	case token.StringVal:
		pv.printf("\"%s\"", s.Value.Lexeme)
	case token.CharVal:
		pv.printf("'%s'", s.Value.Lexeme)
	default:
		pv.printf("%s", s.Value.Lexeme)
	}
	return nil
}

func (pv *PrintVisitor) VisitNewRValue(n *ast.NewRValue) error {
	pv.printf("new %s", n.Type.Lexeme)

	if n.ArrayExpr != nil {
		pv.printf("[")
		if err := n.ArrayExpr.Accept(pv); err != nil {
			return err
		}
		pv.printf("]")
	}

	return nil
}

func (pv *PrintVisitor) VisitVarRValue(v *ast.VarRValue) error {
	if len(v.Path) == 1 {
		pv.printf("%s", v.Path[0].VarName.Lexeme)
		if v.Path[0].ArrayExpr != nil {
			pv.printf("[")
			if err := v.Path[0].ArrayExpr.Accept(pv); err != nil {
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
			}
		}
	}

	return nil
}
