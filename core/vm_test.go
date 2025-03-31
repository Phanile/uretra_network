package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVM_NewVM(t *testing.T) {
	data := []byte{10, 0x01, 20, 0x01, 0x02} // 10 PushInt 20 PushInt Add
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	result := vm.stack.Pop()
	assert.Equal(t, byte(30), result) // stack: 10 20 30
}

func TestVM_Sub(t *testing.T) {
	data := []byte{20, 0x01, 10, 0x01, 0x05} // 20 PushInt 10 PushInt Sub
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	result := vm.stack.Pop()
	assert.Equal(t, byte(10), result) // stack: 20 10 10
}

func TestVM_NewStack(t *testing.T) {
	s := NewStack(8)

	s.Push(1)
	s.Push(4)

	v := s.Pop()

	assert.Equal(t, v, 1)

	v2 := s.Pop()
	assert.Equal(t, v2, 4)
}

func TestVM_PushBytes(t *testing.T) {
	//105 116 32 119 111 114 107 115
	//i   t      w   o   r   k   s
	data := []byte{105, 0x03, 116, 0x03, 32, 0x03, 119, 0x03, 111, 0x03, 114, 0x03, 107, 0x03, 115, 0x03}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	var b []byte
	for i := 0; i < 8; i++ {
		b = append(b, vm.stack.Pop().(byte))
	}
	assert.Equal(t, string(b), "it works")
}

func TestVM_Pack(t *testing.T) {
	//8, PushInt, i, PushBytes, t, PushBytes, space, PushBytes, w, PushBytes, o, PushBytes, r, PushBytes, k, PushBytes, s, PushBytes, Pack
	data := []byte{8, 0x01, 105, 0x03, 116, 0x03, 32, 0x03, 119, 0x03, 111, 0x03, 114, 0x03, 107, 0x03, 115, 0x03, 0x04}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	assert.Equal(t, string(vm.stack.Pop().([]byte)), "it works")
}
