package vm

type OpCode string

const (
	// consts/vars
	PUSH  OpCode = "push"  // [op] push v onto stack
	POP   OpCode = "pop"   // pop val off stack
	LOAD  OpCode = "load"  // [op] push val @ mem addr v onto stack
	STORE OpCode = "store" // [op] pop x, store x @ mem addr v

	// arithmetic
	ADD OpCode = "add" // pop x, y off stack, push (y + x) onto stack
	SUB OpCode = "sub" // pop x, y off stack, push (y - x) onto stack
	MUL OpCode = "mul" // pop x, y off stack, push (y * x) onto stack
	DIV OpCode = "div" //pop x, y off stack, push (y / x) onto stack

	// logical
	AND OpCode = "and" // pop bools x, y, push (y and x)
	OR  OpCode = "or"  // pop bools x, y, push (y or x)
	NOT OpCode = "not" // pop bool x, push (not x)

	// comparators
	CMPLT OpCode = "cmplt" // pop x, y off stack, push (y < x)
	CMPLE OpCode = "cmple" // pop x, y off stack, push (y <= x)
	CMPGT OpCode = "cmpgt" // pop x, y off stack, push (y > x)
	CMPGE OpCode = "cmpge" // pop x, y off stack, push (y >= x)
	CMPEQ OpCode = "cmpeq" // pop x, y off stack, push (y == x)
	CMPNE OpCode = "cmpne" // pop x, y off stack, push (y != x)

	// jump
	JMP  OpCode = "jmp"  // [op] jump to given instr v
	JMPF OpCode = "jmpf" // [op] pop x, if x is F jump to instr v
	CALL OpCode = "call" // [op] call fxn v (pop, push args)
	RET  OpCode = "ret"  // return from curr fxn

	// built-ins
	WRITE  OpCode = "write"  // pop x, write to stdout
	READ   OpCode = "read"   // stdin, push onto stack
	SLEN   OpCode = "slen"   // pop str x, push x.size()
	ALEN   OpCode = "alen"   // pop arr x, push x.size()
	GETC   OpCode = "getc"   // pop str x, pop int y, push x[y]
	TOINT  OpCode = "toint"  // pop x, push x as int
	TODBL  OpCode = "todbl"  // pop x, push x as double
	TOSTR  OpCode = "tostr"  // pop x, push x as str
	CONCAT OpCode = "concat" // pop x, y, push y + x (str concat)

	// heap
	ALLOCS OpCode = "allocs" // alloc struct, push oid x
	ALLOCA OpCode = "alloca" // pop x, y, alloc arr w/ y, x valus, push oid
	ADDF   OpCode = "addf"   // [op] pop x, add field v to obj(x)
	SETF   OpCode = "setf"   // [op] pop x, y, set obj(y).v = x
	GETF   OpCode = "getf"   // [op] pop x, push val of obj(x).v
	SETI   OpCode = "seti"   // pop x, y, z, set arr obj(z)[y] = x
	GETI   OpCode = "geti"   // pop x, y, push arr obj(y)[x] value

	// special
	DUP OpCode = "dup" // pop x, push x, push x
	NOP OpCode = "nop" // no effect, jump code segments
)
