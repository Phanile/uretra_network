package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVM_NewVM(t *testing.T) {
	data := []byte{10, 0x01, 20, 0x01, 0x02} // 10 Push 20 Push Add
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	assert.Equal(t, vm.stack[vm.sp], byte(30)) // stack: 10 20 30
}
