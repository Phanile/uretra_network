package core

import (
	"encoding/binary"
	"fmt"
)

type Instruction byte

const (
	PushInt   Instruction = 0x01 //1
	Add       Instruction = 0x02 //2
	PushBytes Instruction = 0x03 //3
	Pack      Instruction = 0x04 //4
	Sub       Instruction = 0x05 //5
	Store     Instruction = 0x06 //6
)

type VM struct {
	data  []byte
	ip    int //instruction pointer
	stack *Stack
	state *State
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

	o := s.data[0]
	s.data = append(s.data[:0], s.data[1:]...)
	s.sp--

	return o
}

func NewVM(data []byte, state *State) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: NewStack(128),
		state: state,
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

	case Store:
		var (
			key             = vm.stack.Pop().([]byte)
			value           = vm.stack.Pop()
			serializedValue []byte
		)

		switch v := value.(type) {
		case uint8:
			serializedValue = serializeInt64(int64(v))
		default:
			panic("unknown type")
		}

		_ = vm.state.Put(key, serializedValue)
	default:
		break
	}

	return nil
}

func serializeInt64(data int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(data))

	return buf
}

func deserializeInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}
