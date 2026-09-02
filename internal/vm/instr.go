package vm

import "fmt"

type Value interface{}

type Instr struct {
	Opcode  OpCode
	Operand Value
	Comment *string
}

func (i *Instr) Instr(opcode OpCode, operand *Value) {
	i.Opcode = opcode
	if operand != nil {
		i.Operand = *operand
	}
}

func NewInstr(opcode OpCode, operand Value) Instr {
	return Instr{Opcode: opcode, Operand: operand}
}

func (i *Instr) PUSH(value *Value) Instr { return NewInstr(PUSH, *value) }
func (i *Instr) POP() Instr              { return NewInstr(POP, nil) }
func (i *Instr) LOAD(addr int) Instr     { return NewInstr(LOAD, addr) }
func (i *Instr) STORE(addr int) Instr    { return NewInstr(STORE, addr) }

func (i *Instr) ADD() Instr { return NewInstr(ADD, nil) }
func (i *Instr) SUB() Instr { return NewInstr(SUB, nil) }
func (i *Instr) MUL() Instr { return NewInstr(MUL, nil) }
func (i *Instr) DIV() Instr { return NewInstr(DIV, nil) }

func (i *Instr) AND() Instr { return NewInstr(AND, nil) }
func (i *Instr) OR() Instr  { return NewInstr(OR, nil) }
func (i *Instr) NOT() Instr { return NewInstr(NOT, nil) }

func (i *Instr) CMPLT() Instr { return NewInstr(CMPLT, nil) }
func (i *Instr) CMPLE() Instr { return NewInstr(CMPLE, nil) }
func (i *Instr) CMPGT() Instr { return NewInstr(CMPGT, nil) }
func (i *Instr) CMPGE() Instr { return NewInstr(CMPGE, nil) }
func (i *Instr) CMPEQ() Instr { return NewInstr(CMPEQ, nil) }
func (i *Instr) CMPNE() Instr { return NewInstr(CMPNE, nil) }

func (i *Instr) JMP(index int) Instr    { return NewInstr(JMP, index) }
func (i *Instr) JMPF(index int) Instr   { return NewInstr(JMPF, index) }
func (i *Instr) CALL(fun *string) Instr { return NewInstr(CALL, *fun) }
func (i *Instr) RET() Instr             { return NewInstr(RET, nil) }

func (i *Instr) WRITE() Instr  { return NewInstr(WRITE, nil) }
func (i *Instr) READ() Instr   { return NewInstr(READ, nil) }
func (i *Instr) SLEN() Instr   { return NewInstr(SLEN, nil) }
func (i *Instr) ALEN() Instr   { return NewInstr(ALEN, nil) }
func (i *Instr) GETC() Instr   { return NewInstr(GETC, nil) }
func (i *Instr) TOINT() Instr  { return NewInstr(TOINT, nil) }
func (i *Instr) TODBL() Instr  { return NewInstr(TODBL, nil) }
func (i *Instr) TOSTR() Instr  { return NewInstr(TOSTR, nil) }
func (i *Instr) CONCAT() Instr { return NewInstr(CONCAT, nil) }

func (i *Instr) ALLOCS() Instr            { return NewInstr(ALLOCS, nil) }
func (i *Instr) ALLOCA() Instr            { return NewInstr(ALLOCA, nil) }
func (i *Instr) ADDF(field *string) Instr { return NewInstr(ADDF, *field) }
func (i *Instr) SETF(field *string) Instr { return NewInstr(SETF, *field) }
func (i *Instr) GETF(field *string) Instr { return NewInstr(GETF, *field) }

func (i *Instr) SETI() Instr { return NewInstr(SETI, nil) }
func (i *Instr) GETI() Instr { return NewInstr(GETI, nil) }

func (i *Instr) DUP() Instr { return NewInstr(DUP, nil) }
func (i *Instr) NOP() Instr { return NewInstr(NOP, nil) }

func (i *Instr) String(instr *Instr) string {
	if instr == nil {
		return "<nil>"
	}
	if instr.Operand == nil {
		return string(instr.Opcode) + "()"
	}
	return fmt.Sprintf("%s(%v)", instr.Opcode, instr.Operand)
}

func ValueString(value Value) string {
	if value == nil {
		return "null"
	}
	return fmt.Sprint(value)
}
