package core

type Instruction byte

const (
	Push Instruction = 0x01 //1
	Add  Instruction = 0x02 //2
)

type VM struct {
	data  []byte
	ip    int //instruction pointer
	stack []byte
	sp    int //stack pointer
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: make([]byte, 1024),
		sp:    -1,
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
	case Push:
		vm.pushStack(vm.data[vm.ip-1])
	case Add:
		a := vm.stack[0]
		b := vm.stack[1]
		c := a + b
		vm.pushStack(c)
	default:
		break
	}

	return nil
}

func (vm *VM) pushStack(b byte) {
	vm.sp++
	vm.stack[vm.sp] = b
}
