package core

import "fmt"

type Instruction byte

const (
	PushInt   Instruction = 0x01 //1
	Add       Instruction = 0x02 //2
	PushBytes Instruction = 0x03 //3
	Pack      Instruction = 0x04 //4
	Sub       Instruction = 0x05 //5
)

type VM struct {
	data  []byte
	ip    int //instruction pointer
	stack *Stack
}

type Stack struct {
	data []any
	sp   int // stack pointer
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

func (s *Stack) Push(o any) {
	s.data[s.sp] = o
	s.sp++
}

func (s *Stack) Pop() any {
	if s.sp == 0 {
		return fmt.Errorf("no data in the stack")
	}
	s.sp--
	o := s.data[0]
	s.data = append(s.data[:0], s.data[1:]...)

	return o
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: NewStack(128),
	}
}

func (vm *VM) Run() error {
	for {
		instruction := Instruction(vm.data[vm.ip])

		if err := vm.Execute(instruction); err != nil {
			return err
		}

		vm.ip++

		if vm.ip > len(vm.data)-1 {
			break
		}
	}
	return nil
}

func (vm *VM) Execute(instr Instruction) error {
	switch instr {
	case PushInt:
		vm.stack.Push(vm.data[vm.ip-1])
	case PushBytes:
		vm.stack.Push(vm.data[vm.ip-1])
	case Add:
		a := vm.stack.Pop().(uint8)
		b := vm.stack.Pop().(uint8)
		c := a + b
		vm.stack.Push(c)
	case Sub:
		a := vm.stack.Pop().(uint8)
		b := vm.stack.Pop().(uint8)
		c := a - b
		vm.stack.Push(c)
	case Pack:
		n := vm.stack.Pop().(uint8)
		b := make([]byte, n)

		for i := 0; i < int(n); i++ {
			b[i] = vm.stack.Pop().(byte)
		}

		vm.stack.Push(b)
	default:
		break
	}

	return nil
}
