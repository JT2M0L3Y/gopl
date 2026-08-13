package ast

import (
	"fmt"
	"gopl/internal/lexer"
	"gopl/internal/token"
)

// Parser performs recursive descent parsing on tokens
type Parser struct {
	lexer     *lexer.Lexer
	currToken token.Token
}

// New creates a new parser from a lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.advance()
	return p
}

// Parse parses a complete program
func (p *Parser) Parse() (*Program, error) {
	prog := &Program{}

	for !p.match(token.EOF) {
		if p.match(token.Struct) {
			if err := p.structDef(prog); err != nil {
				return nil, err
			}
		} else {
			if err := p.funDef(prog); err != nil {
				return nil, err
			}
		}
	}

	if err := p.eat(token.EOF, "expecting end-of-file"); err != nil {
		return nil, err
	}

	return prog, nil
}

// Private helper methods

func (p *Parser) advance() {
	p.currToken = p.lexer.Next()
}

func (p *Parser) eat(k token.Kind, msg string) error {
	if !p.match(k) {
		return p.error(msg)
	}
	p.advance()
	return nil
}

func (p *Parser) match(k token.Kind) bool {
	return p.currToken.Kind == k
}

func (p *Parser) matchAny(kinds []token.Kind) bool {
	for _, k := range kinds {
		if p.match(k) {
			return true
		}
	}
	return false
	// ! Go slices can replace this
	// return slices.ContainsFunc(kinds, p.match)
}

func (p *Parser) error(msg string) error {
	return fmt.Errorf("%s found '%s' at line %d, column %d",
		msg, p.currToken.Lexeme, p.currToken.Line, p.currToken.Col)
}

func (p *Parser) binOp() bool {
	return p.matchAny([]token.Kind{
		token.Plus, token.Minus, token.Times, token.Divide,
		token.And, token.Or, token.Equal, token.Less,
		token.Greater, token.LessEq, token.GreaterEq, token.NotEqual,
	})
}

func (p *Parser) baseType() bool {
	return p.matchAny([]token.Kind{
		token.IntType, token.FloatType, token.BoolType,
		token.CharType, token.StringType,
	})
}

func (p *Parser) baseRValue() bool {
	return p.matchAny([]token.Kind{
		token.IntVal, token.FloatVal, token.BoolVal,
		token.CharVal, token.StringVal,
	})
}

// Parsing functions for top-level constructs

func (p *Parser) structDef(prog *Program) error {
	if err := p.eat(token.Struct, "expecting 'struct'"); err != nil {
		return err
	}

	s := StructDef{
		StructName: p.currToken,
	}

	if err := p.eat(token.Ident, "expecting identifier"); err != nil {
		return err
	}

	if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
		return err
	}

	if err := p.fields(&s); err != nil {
		return err
	}

	if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
		return err
	}

	prog.StructDefs = append(prog.StructDefs, s)
	return nil
}

func (p *Parser) fields(s *StructDef) error {
	if p.baseType() || p.matchAny([]token.Kind{token.Ident, token.Array}) {
		v := VarDef{}
		if err := p.dataType(&v); err != nil {
			return err
		}
		v.VarName = p.currToken
		if err := p.eat(token.Ident, "expecting identifier"); err != nil {
			return err
		}
		s.Fields = append(s.Fields, v)

		for p.match(token.Comma) {
			p.advance()
			v := VarDef{}
			if err := p.dataType(&v); err != nil {
				return err
			}
			v.VarName = p.currToken
			if err := p.eat(token.Ident, "expecting identifier"); err != nil {
				return err
			}
			s.Fields = append(s.Fields, v)
		}
	}
	return nil
}

func (p *Parser) funDef(prog *Program) error {
	fun := FuncDef{}

	v := VarDef{}
	if err := p.dataType(&v); err != nil {
		return err
	}
	fun.ReturnType = v.DataType

	fun.FunName = p.currToken
	if err := p.eat(token.Ident, "expecting function name"); err != nil {
		return err
	}

	if err := p.eat(token.LParen, "expecting '('"); err != nil {
		return err
	}

	if err := p.params(&fun); err != nil {
		return err
	}

	if err := p.eat(token.RParen, "expecting ')'"); err != nil {
		return err
	}

	if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
		return err
	}

	var stmts []Stmt
	if err := p.stmt(&stmts); err != nil {
		return err
	}
	fun.Stmts = stmts

	if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
		return err
	}

	prog.FunDefs = append(prog.FunDefs, fun)
	return nil
}

func (p *Parser) params(f *FuncDef) error {
	for !p.match(token.RParen) {
		v := VarDef{}
		if err := p.dataType(&v); err != nil {
			return err
		}

		v.VarName = p.currToken
		if err := p.eat(token.Ident, "expecting identifier"); err != nil {
			return err
		}

		f.Params = append(f.Params, v)

		if p.match(token.Comma) {
			p.advance()
			if p.match(token.RParen) {
				return p.error("expecting another parameter")
			}
		}
	}
	return nil
}

func (p *Parser) dataType(v *VarDef) error {
	if p.match(token.VoidType) {
		v.DataType.TypeNames = append(v.DataType.TypeNames, p.currToken.Lexeme)
		p.advance()
	} else if p.baseType() {
		v.DataType.TypeNames = append(v.DataType.TypeNames, p.currToken.Lexeme)
		p.advance()
	} else if p.match(token.Ident) {
		v.DataType.TypeNames = append(v.DataType.TypeNames, p.currToken.Lexeme)
		p.advance()
	} else if p.match(token.Array) {
		v.DataType.IsArray = true
		p.advance()
		if p.baseType() || p.match(token.Ident) {
			v.DataType.TypeNames = append(v.DataType.TypeNames, p.currToken.Lexeme)
			p.advance()
		}
	} else if p.match(token.Dict) {
		v.DataType.IsDict = true
		p.advance()
		// parse key type (typically string)
		if p.match(token.StringType) {
			v.DataType.TypeNames = append(v.DataType.TypeNames, p.currToken.Lexeme)
			p.advance()
		}
		// parse value type
		if p.baseType() {
			v.DataType.TypeNames = append(v.DataType.TypeNames, p.currToken.Lexeme)
			p.advance()
		}
	} else {
		return p.error("expecting type")
	}
	return nil
}

// Statement parsing

func (p *Parser) stmt(stmts *[]Stmt) error {
	for !p.match(token.RBrace) {
		if p.match(token.If) {
			iStmt := &IfStmt{}
			if err := p.ifStmt(iStmt); err != nil {
				return err
			}
			*stmts = append(*stmts, iStmt)
		} else if p.match(token.For) {
			fStmt := &ForStmt{}
			if err := p.forStmt(fStmt); err != nil {
				return err
			}
			*stmts = append(*stmts, fStmt)
		} else if p.match(token.While) {
			wStmt := &WhileStmt{}
			if err := p.whileStmt(wStmt); err != nil {
				return err
			}
			*stmts = append(*stmts, wStmt)
		} else if p.match(token.Return) {
			rStmt := &ReturnStmt{}
			if err := p.retStmt(rStmt); err != nil {
				return err
			}
			*stmts = append(*stmts, rStmt)
		} else if p.matchAny([]token.Kind{token.Array, token.Dict}) || p.baseType() {
			vDecl := &VarDeclStmt{}
			if err := p.dataType(&vDecl.VarDef); err != nil {
				return err
			}
			if err := p.vdeclStmt(vDecl); err != nil {
				return err
			}
			*stmts = append(*stmts, vDecl)
		} else if p.match(token.Ident) {
			idToken := p.currToken
			p.advance()

			if p.matchAny([]token.Kind{token.Assign, token.Dot}) {
				aStmt := &AssignStmt{}
				v := VarRef{VarName: idToken}
				aStmt.LValue = append(aStmt.LValue, v)
				if err := p.assignStmt(aStmt, idToken); err != nil {
					return err
				}
				*stmts = append(*stmts, aStmt)
			} else if p.match(token.LBracket) {
				p.advance()
				aStmt := &AssignStmt{}
				e := &Expr{}
				if err := p.expr(e); err != nil {
					return err
				}
				v := VarRef{VarName: idToken}
				if e.First != nil && e.First.FirstToken().Kind == token.StringVal {
					v.DictExpr = e
				} else {
					v.ArrayExpr = e
				}
				aStmt.LValue = append(aStmt.LValue, v)
				if err := p.eat(token.RBracket, "expecting ']'"); err != nil {
					return err
				}
				if err := p.assignStmt(aStmt, idToken); err != nil {
					return err
				}
				*stmts = append(*stmts, aStmt)
			} else if p.match(token.LParen) {
				cExpr := &CallExpr{FunName: idToken}
				if err := p.callExpr(cExpr); err != nil {
					return err
				}
				*stmts = append(*stmts, cExpr)
			} else {
				// implicit variable declaration with type inferred from context
				vDecl := &VarDeclStmt{}
				vDecl.VarDef.VarName = idToken
				vDecl.VarDef.DataType.TypeNames = append(vDecl.VarDef.DataType.TypeNames, idToken.Lexeme)
				if err := p.vdeclStmt(vDecl); err != nil {
					return err
				}
				*stmts = append(*stmts, vDecl)
			}
		} else {
			return p.error("expecting statement")
		}
	}
	return nil
}

func (p *Parser) vdeclStmt(vDecl *VarDeclStmt) error {
	if !p.match(token.Assign) {
		vDecl.VarDef.VarName = p.currToken
		if err := p.eat(token.Ident, "expecting identifier"); err != nil {
			return err
		}
	}
	if err := p.eat(token.Assign, "expecting '='"); err != nil {
		return err
	}
	e := &Expr{}
	if err := p.expr(e); err != nil {
		return err
	}
	vDecl.Expr = *e
	return nil
}

func (p *Parser) assignStmt(aStmt *AssignStmt, t token.Token) error {
	for !p.match(token.Assign) {
		v := VarRef{}
		if err := p.lvalue(&v, t); err != nil {
			return err
		}
		aStmt.LValue = append(aStmt.LValue, v)
	}
	if err := p.eat(token.Assign, "expecting '='"); err != nil {
		return err
	}
	e := &Expr{}
	if err := p.expr(e); err != nil {
		return err
	}
	aStmt.Expr = *e
	return nil
}

func (p *Parser) lvalue(v *VarRef, t token.Token) error {
	if p.match(token.Ident) {
		v.VarName = p.currToken
		p.advance()
	}

	if p.match(token.Dot) {
		p.advance()
		v.VarName = p.currToken
		if err := p.eat(token.Ident, "expecting identifier"); err != nil {
			return err
		}
	}

	if p.match(token.LBracket) {
		p.advance()
		e := &Expr{}
		if err := p.expr(e); err != nil {
			return err
		}
		v.ArrayExpr = e
		if err := p.eat(token.RBracket, "expecting ']'"); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) ifStmt(iStmt *IfStmt) error {
	if err := p.eat(token.If, "expecting 'if'"); err != nil {
		return err
	}
	if err := p.eat(token.LParen, "expecting '('"); err != nil {
		return err
	}

	e := &Expr{}
	if err := p.expr(e); err != nil {
		return err
	}
	iStmt.IfPart.Condition = *e

	if err := p.eat(token.RParen, "expecting ')'"); err != nil {
		return err
	}

	if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
		return err
	}

	var stmts []Stmt
	if err := p.stmt(&stmts); err != nil {
		return err
	}
	iStmt.IfPart.Stmts = stmts

	if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
		return err
	}

	return p.ifStmtT(iStmt)
}

func (p *Parser) ifStmtT(i *IfStmt) error {
	if p.match(token.ElseIf) {
		bi := BasicIf{}

		p.advance()

		e := &Expr{}
		if err := p.expr(e); err != nil {
			return err
		}
		bi.Condition = *e

		if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
			return err
		}

		var stmts []Stmt
		if err := p.stmt(&stmts); err != nil {
			return err
		}
		bi.Stmts = stmts

		if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
			return err
		}
		i.ElseIfs = append(i.ElseIfs, bi)

		return p.ifStmtT(i)
	} else if p.match(token.Else) {
		p.advance()
		if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
			return err
		}

		var stmts []Stmt
		if err := p.stmt(&stmts); err != nil {
			return err
		}
		i.ElseStmts = stmts

		if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) whileStmt(wStmt *WhileStmt) error {
	if err := p.eat(token.While, "expecting 'while'"); err != nil {
		return err
	}

	e := &Expr{}
	if err := p.expr(e); err != nil {
		return err
	}
	wStmt.Condition = *e

	if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
		return err
	}

	var stmts []Stmt
	if err := p.stmt(&stmts); err != nil {
		return err
	}
	wStmt.Stmts = stmts

	if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
		return err
	}
	return nil
}

func (p *Parser) forStmt(fStmt *ForStmt) error {
	if err := p.eat(token.For, "expecting 'for'"); err != nil {
		return err
	}

	if err := p.eat(token.LParen, "expecting '('"); err != nil {
		return err
	}

	vDecl := &VarDeclStmt{}
	if err := p.dataType(&vDecl.VarDef); err != nil {
		return err
	}
	if err := p.vdeclStmt(vDecl); err != nil {
		return err
	}
	fStmt.VarDecl = *vDecl
	if err := p.eat(token.Semicolon, "expecting ';'"); err != nil {
		return err
	}

	e := &Expr{}
	if err := p.expr(e); err != nil {
		return err
	}
	fStmt.Condition = *e

	if err := p.eat(token.Semicolon, "expecting ';'"); err != nil {
		return err
	}

	aStmt := &AssignStmt{}
	if err := p.assignStmt(aStmt, fStmt.VarDecl.VarDef.VarName); err != nil {
		return err
	}
	fStmt.Assign = *aStmt

	if err := p.eat(token.RParen, "expecting ')'"); err != nil {
		return err
	}

	if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
		return err
	}

	var stmts []Stmt
	if err := p.stmt(&stmts); err != nil {
		return err
	}
	fStmt.Stmts = stmts

	if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
		return err
	}
	return nil
}

func (p *Parser) callExpr(cExpr *CallExpr) error {
	if err := p.eat(token.LParen, "expecting '('"); err != nil {
		return err
	}

	if p.baseRValue() || p.matchAny([]token.Kind{token.NullVal, token.New, token.Ident}) {
		e := &Expr{}
		if err := p.expr(e); err != nil {
			return err
		}
		cExpr.Args = append(cExpr.Args, *e)

		for p.match(token.Comma) {
			p.advance()
			e := &Expr{}
			if err := p.expr(e); err != nil {
				return err
			}
			cExpr.Args = append(cExpr.Args, *e)
		}
	}

	if err := p.eat(token.RParen, "expecting ')'"); err != nil {
		return err
	}
	return nil
}

func (p *Parser) retStmt(rStmt *ReturnStmt) error {
	if err := p.eat(token.Return, "expecting 'return'"); err != nil {
		return err
	}
	e := &Expr{}
	if err := p.expr(e); err != nil {
		return err
	}
	rStmt.Expr = *e
	return nil
}

// Expression parsing

func (p *Parser) expr(e *Expr) error {
	for p.match(token.Not) {
		e.Negated = true
		p.advance()
	}

	if p.match(token.LParen) {
		p.advance()
		cTerm := &ComplexTerm{}
		inner := &Expr{}
		if err := p.expr(inner); err != nil {
			return err
		}
		cTerm.Expr = *inner
		e.First = cTerm
		if err := p.eat(token.RParen, "expecting ')'"); err != nil {
			return err
		}
	} else if p.baseRValue() || p.baseType() || p.matchAny([]token.Kind{token.NullVal, token.New, token.Ident}) {
		sTerm := &SimpleTerm{}
		if err := p.rvalue(sTerm); err != nil {
			return err
		}
		e.First = sTerm
	} else {
		return p.error("expecting expression start")
	}

	if p.binOp() || p.match(token.Comma) {
		op := p.currToken
		e.Op = &op
		p.advance()
		restExpr := &Expr{}
		if err := p.expr(restExpr); err != nil {
			return err
		}
		e.Rest = restExpr
	}

	return nil
}

func (p *Parser) rvalue(sTerm *SimpleTerm) error {
	if p.matchAny([]token.Kind{token.NullVal}) || p.baseRValue() || p.baseType() {
		sRVal := &SimpleRValue{Value: p.currToken}
		p.advance()
		sTerm.RValue = sRVal
	} else if p.match(token.New) {
		p.advance()
		newRVal := &NewRValue{}
		if err := p.newRValue(newRVal); err != nil {
			return err
		}
		sTerm.RValue = newRVal
	} else if p.match(token.Ident) {
		tmp := p.currToken
		p.advance()

		if p.match(token.LParen) {
			cExpr := &CallExpr{FunName: tmp}
			if err := p.callExpr(cExpr); err != nil {
				return err
			}
			sTerm.RValue = cExpr
		} else {
			v := VarRef{VarName: tmp}
			vRVal := &VarRValue{}
			vRVal.Path = append(vRVal.Path, v)
			if err := p.varRValue(vRVal); err != nil {
				return err
			}
			sTerm.RValue = vRVal
		}
	} else {
		return p.error("expecting rvalue")
	}
	return nil
}

func (p *Parser) newRValue(nRVal *NewRValue) error {
	if p.match(token.Ident) {
		nRVal.Type = p.currToken
		p.advance()

		if p.match(token.LBracket) {
			p.advance()
			e := &Expr{}
			if err := p.expr(e); err != nil {
				return err
			}
			nRVal.ArrayExpr = e
			if err := p.eat(token.RBracket, "expecting ']'"); err != nil {
				return err
			}
		}
	} else if p.match(token.Dict) {
		nRVal.Type = p.currToken
		p.advance()
		if err := p.eat(token.LBrace, "expecting '{'"); err != nil {
			return err
		}
		e := &Expr{}
		if err := p.expr(e); err != nil {
			return err
		}
		nRVal.DictExpr = e
		if err := p.eat(token.RBrace, "expecting '}'"); err != nil {
			return err
		}
	} else if p.baseType() {
		nRVal.Type = p.currToken
		p.advance()
		if err := p.eat(token.LBracket, "expecting '['"); err != nil {
			return err
		}
		e := &Expr{}
		if err := p.expr(e); err != nil {
			return err
		}
		nRVal.ArrayExpr = e
		if err := p.eat(token.RBracket, "expecting ']'"); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) varRValue(vRVal *VarRValue) error {
	for p.matchAny([]token.Kind{token.Dot, token.LBracket}) {
		if p.match(token.Dot) {
			p.advance()
			v := VarRef{VarName: p.currToken}
			vRVal.Path = append(vRVal.Path, v)
			if err := p.eat(token.Ident, "expecting identifier"); err != nil {
				return err
			}
		} else if p.match(token.LBracket) {
			p.advance()
			e := &Expr{}
			if err := p.expr(e); err != nil {
				return err
			}
			if e.First != nil && e.First.FirstToken().Kind == token.StringVal {
				vRVal.Path[len(vRVal.Path)-1].DictExpr = e
			} else {
				vRVal.Path[len(vRVal.Path)-1].ArrayExpr = e
			}
			if err := p.eat(token.RBracket, "expecting ']'"); err != nil {
				return err
			}
		}
	}
	return nil
}
