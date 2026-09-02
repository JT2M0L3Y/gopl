package vm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopl/internal/utils"
)

type VM struct {
	structHeap map[int]map[string]Value
	arrayHeap  map[int][]Value
	nextObjID  int
	frameInfo  map[string]FrameInfo
	callStack  *utils.Stack[*Frame]
	input      io.Reader
	output     io.Writer
}

func New() *VM {
	return &VM{
		structHeap: map[int]map[string]Value{},
		arrayHeap:  map[int][]Value{},
		nextObjID:  2023,
		frameInfo:  map[string]FrameInfo{},
		callStack:  utils.New[*Frame](),
		input:      os.Stdin,
		output:     os.Stdout,
	}
}

func (vm *VM) SetInput(input io.Reader) {
	if input != nil {
		vm.input = input
	}
}

func (vm *VM) SetOutput(output io.Writer) {
	if output != nil {
		vm.output = output
	}
}

func (vm *VM) Add(frame *FrameInfo) {
	if vm.frameInfo == nil {
		vm.frameInfo = map[string]FrameInfo{}
	}
	if vm.structHeap == nil {
		vm.structHeap = map[int]map[string]Value{}
	}
	if vm.arrayHeap == nil {
		vm.arrayHeap = map[int][]Value{}
	}
	if vm.nextObjID == 0 {
		vm.nextObjID = 2023
	}
	if vm.callStack == nil {
		vm.callStack = utils.New[*Frame]()
	}
	vm.frameInfo[frame.Name] = *frame
}

func (vm *VM) Execute(debug bool) error {
	if vm.callStack == nil {
		vm.callStack = utils.New[*Frame]()
	}
	info, ok := vm.frameInfo["main"]
	if !ok {
		return fmt.Errorf("VM error: no 'main' function")
	}
	frame := &Frame{Info: info}
	if err := vm.callStack.Push(frame); err != nil {
		return err
	}
	reader := bufio.NewReader(vm.input)
	for !vm.callStack.IsEmpty() {
		if frame.ProgCount >= len(frame.Info.Instructions) {
			vm.callStack.Pop()
			if vm.callStack.IsEmpty() {
				break
			}
			frame, _ = vm.callStack.Peek()
			continue
		}
		pc := frame.ProgCount
		instr := frame.Info.Instructions[pc]
		frame.ProgCount++
		if debug {
			fmt.Fprintf(os.Stderr, "%s[%d] %s\n", frame.Info.Name, pc, (&instr).String(&instr))
		}
		if err := vm.executeInstr(&frame, instr, reader); err != nil {
			return fmt.Errorf("%v (in %s at %d)", err, frame.Info.Name, pc)
		}
	}
	return nil
}

func (vm *VM) Run(debug *bool) { enabled := debug != nil && *debug; _ = vm.Execute(enabled) }

func pop(frame *Frame) (Value, error) {
	if len(frame.OpStack) == 0 {
		return nil, fmt.Errorf("operand stack underflow")
	}
	v := frame.OpStack[len(frame.OpStack)-1]
	frame.OpStack = frame.OpStack[:len(frame.OpStack)-1]
	return v, nil
}

func push(frame *Frame, value Value) { frame.OpStack = append(frame.OpStack, value) }

func intValue(v Value) (int, error) {
	n, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("expected integer, got %s", ValueString(v))
	}
	return n, nil
}

func boolValue(v Value) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("expected boolean, got %s", ValueString(v))
	}
	return b, nil
}

func number(v Value) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case float64:
		return n, true
	}
	return 0, false
}

func binary(frame *Frame, fn func(Value, Value) (Value, error)) error {
	x, err := pop(frame)
	if err != nil {
		return err
	}
	y, err := pop(frame)
	if err != nil {
		return err
	}
	v, err := fn(y, x)
	if err != nil {
		return err
	}
	push(frame, v)
	return nil
}

func (vm *VM) executeInstr(frame **Frame, instr Instr, reader *bufio.Reader) error {
	f := *frame
	switch instr.Opcode {
	case PUSH:
		push(f, instr.Operand)
	case POP:
		_, err := pop(f)
		return err
	case LOAD:
		i, err := intValue(instr.Operand)
		if err != nil {
			return err
		}
		if i < 0 || i >= len(f.Variables) {
			return fmt.Errorf("invalid variable address %d", i)
		}
		push(f, f.Variables[i])
	case STORE:
		i, err := intValue(instr.Operand)
		if err != nil {
			return err
		}
		v, err := pop(f)
		if err != nil {
			return err
		}
		for len(f.Variables) <= i {
			f.Variables = append(f.Variables, nil)
		}
		f.Variables[i] = v
	case ADD:
		return binary(f, func(y, x Value) (Value, error) { return arithmetic(y, x, '+') })
	case SUB:
		return binary(f, func(y, x Value) (Value, error) { return arithmetic(y, x, '-') })
	case MUL:
		return binary(f, func(y, x Value) (Value, error) { return arithmetic(y, x, '*') })
	case DIV:
		return binary(f, func(y, x Value) (Value, error) { return arithmetic(y, x, '/') })
	case AND, OR:
		return binary(f, func(y, x Value) (Value, error) {
			a, e := boolValue(y)
			if e != nil {
				return nil, e
			}
			b, e := boolValue(x)
			if e != nil {
				return nil, e
			}
			if instr.Opcode == AND {
				return a && b, nil
			}
			return a || b, nil
		})
	case NOT:
		v, err := pop(f)
		if err != nil {
			return err
		}
		b, err := boolValue(v)
		if err != nil {
			return err
		}
		push(f, !b)
	case CMPLT, CMPLE, CMPGT, CMPGE:
		return binary(f, func(y, x Value) (Value, error) {
			a, ok := number(y)
			b, ok2 := number(x)
			if ok && ok2 {
				switch instr.Opcode {
				case CMPLT:
					return a < b, nil
				case CMPLE:
					return a <= b, nil
				case CMPGT:
					return a > b, nil
				default:
					return a >= b, nil
				}
			}
			ys, ok := y.(string)
			xs, ok2 := x.(string)
			if !ok || !ok2 {
				return nil, fmt.Errorf("values are not comparable")
			}
			switch instr.Opcode {
			case CMPLT:
				return ys < xs, nil
			case CMPLE:
				return ys <= xs, nil
			case CMPGT:
				return ys > xs, nil
			default:
				return ys >= xs, nil
			}
		})
	case CMPEQ, CMPNE:
		return binary(f, func(y, x Value) (Value, error) {
			equal := fmt.Sprint(y) == fmt.Sprint(x)
			if y == nil || x == nil {
				equal = y == nil && x == nil
			}
			if instr.Opcode == CMPNE {
				equal = !equal
			}
			return equal, nil
		})
	case JMP:
		n, err := intValue(instr.Operand)
		if err != nil {
			return err
		}
		f.ProgCount = n
	case JMPF:
		n, err := intValue(instr.Operand)
		if err != nil {
			return err
		}
		v, err := pop(f)
		if err != nil {
			return err
		}
		b, err := boolValue(v)
		if err != nil {
			return err
		}
		if !b {
			f.ProgCount = n
		}
	case CALL:
		name, ok := instr.Operand.(string)
		if !ok {
			return fmt.Errorf("invalid call operand")
		}
		info, ok := vm.frameInfo[name]
		if !ok {
			return fmt.Errorf("unknown function '%s'", name)
		}
		callee := &Frame{Info: info}
		for i := 0; i < info.ArgCount; i++ {
			v, err := pop(f)
			if err != nil {
				return err
			}
			callee.OpStack = append(callee.OpStack, v)
		}
		if err := vm.callStack.Push(callee); err != nil {
			return err
		}
		*frame = callee
	case RET:
		result, err := pop(f)
		if err != nil {
			return err
		}
		vm.callStack.Pop()
		if vm.callStack.IsEmpty() {
			return nil
		}
		caller, _ := vm.callStack.Peek()
		push(caller, result)
		*frame = caller
	case WRITE:
		v, err := pop(f)
		if err != nil {
			return err
		}
		fmt.Fprint(vm.output, ValueString(v))
	case READ:
		s, err := reader.ReadString('\n')
		if err != nil && len(s) == 0 {
			return err
		}
		push(f, strings.TrimSuffix(strings.TrimSuffix(s, "\n"), "\r"))
	case SLEN:
		v, err := pop(f)
		if err != nil {
			return err
		}
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("slen expects string")
		}
		push(f, len(s))
	case ALEN:
		id, err := pop(f)
		if err != nil {
			return err
		}
		n, err := intValue(id)
		if err != nil {
			return err
		}
		push(f, len(vm.arrayHeap[n]))
	case GETC:
		idx, err := pop(f)
		if err != nil {
			return err
		}
		value, err := pop(f)
		if err != nil {
			return err
		}
		i, err := intValue(idx)
		s, ok := value.(string)
		if err != nil || !ok || i < 0 || i >= len(s) {
			return fmt.Errorf("out-of-bounds string index")
		}
		push(f, string(s[i]))
	case TOINT:
		v, err := pop(f)
		if err != nil {
			return err
		}
		n, err := strconv.Atoi(fmt.Sprint(v))
		if err != nil {
			return err
		}
		push(f, n)
	case TODBL:
		v, err := pop(f)
		if err != nil {
			return err
		}
		n, err := strconv.ParseFloat(fmt.Sprint(v), 64)
		if err != nil {
			return err
		}
		push(f, n)
	case TOSTR:
		v, err := pop(f)
		if err != nil {
			return err
		}
		push(f, ValueString(v))
	case CONCAT:
		return binary(f, func(y, x Value) (Value, error) { return fmt.Sprint(y) + fmt.Sprint(x), nil })
	case ALLOCS:
		id := vm.nextObjID
		vm.nextObjID++
		vm.structHeap[id] = map[string]Value{}
		push(f, id)
	case ALLOCA:
		value, err := pop(f)
		if err != nil {
			return err
		}
		size, err := pop(f)
		if err != nil {
			return err
		}
		n, err := intValue(size)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid array size")
		}
		id := vm.nextObjID
		vm.nextObjID++
		vm.arrayHeap[id] = make([]Value, n)
		for i := range vm.arrayHeap[id] {
			vm.arrayHeap[id][i] = value
		}
		push(f, id)
	case ADDF:
		id, err := pop(f)
		if err != nil {
			return err
		}
		n, err := intValue(id)
		if err != nil {
			return err
		}
		vm.structHeap[n][fmt.Sprint(instr.Operand)] = nil
	case SETF:
		value, err := pop(f)
		if err != nil {
			return err
		}
		id, err := pop(f)
		if err != nil {
			return err
		}
		n, err := intValue(id)
		if err != nil {
			return err
		}
		vm.structHeap[n][fmt.Sprint(instr.Operand)] = value
	case GETF:
		id, err := pop(f)
		if err != nil {
			return err
		}
		n, err := intValue(id)
		if err != nil {
			return err
		}
		push(f, vm.structHeap[n][fmt.Sprint(instr.Operand)])
	case SETI:
		value, err := pop(f)
		if err != nil {
			return err
		}
		index, err := pop(f)
		if err != nil {
			return err
		}
		id, err := pop(f)
		if err != nil {
			return err
		}
		i, e := intValue(index)
		n, e2 := intValue(id)
		if e != nil || e2 != nil || i < 0 || i >= len(vm.arrayHeap[n]) {
			return fmt.Errorf("out-of-bounds array index")
		}
		vm.arrayHeap[n][i] = value
	case GETI:
		index, err := pop(f)
		if err != nil {
			return err
		}
		id, err := pop(f)
		if err != nil {
			return err
		}
		i, e := intValue(index)
		n, e2 := intValue(id)
		if e != nil || e2 != nil || i < 0 || i >= len(vm.arrayHeap[n]) {
			return fmt.Errorf("out-of-bounds array index")
		}
		push(f, vm.arrayHeap[n][i])
	case DUP:
		v, err := pop(f)
		if err != nil {
			return err
		}
		push(f, v)
		push(f, v)
	case NOP:
	default:
		return fmt.Errorf("unsupported operation %s", instr.Opcode)
	}
	return nil
}

func arithmetic(y, x Value, op byte) (Value, error) {
	a, ok := number(y)
	b, ok2 := number(x)
	if !ok || !ok2 {
		return nil, fmt.Errorf("arithmetic expects numbers")
	}
	if op == '/' && b == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	if _, ints := y.(int); ints {
		if _, ok := x.(int); ok && op != '/' {
			switch op {
			case '+':
				return int(a + b), nil
			case '-':
				return int(a - b), nil
			case '*':
				return int(a * b), nil
			}
		}
	}
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	default:
		return a / b, nil
	}
}
