package types

import (
	"fmt"
	"reflect"
)

type List[T any] struct {
	Data []T
}

func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

func (l *List[T]) Get(index int) T {
	if index < 0 || index > len(l.Data)-1 {
		err := fmt.Errorf("index out of range")
		panic(err)
	}

	return l.Data[index]
}

func (l *List[T]) Insert(elem T) {
	l.Data = append(l.Data, elem)
}

func (l *List[T]) Clear() {
	l.Data = []T{}
}

func (l *List[T]) GetIndex(elem T) int {
	for i := 0; i < len(l.Data); i++ {
		if reflect.DeepEqual(l.Data[i], elem) {
			return i
		}
	}
	return -1
}

func (l *List[T]) Remove(elem T) {
	index := l.GetIndex(elem)
	if index == -1 {
		err := fmt.Errorf("element not found")
		panic(err)
	}

	l.Pop(index)
}

func (l *List[T]) Pop(index int) {
	if index < 0 || index > len(l.Data)-1 {
		err := fmt.Errorf("index out of range")
		panic(err)
	}

	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

func (l *List[T]) Contains(elem T) bool {
	for i := 0; i < len(l.Data); i++ {
		if reflect.DeepEqual(l.Data[i], elem) {
			return true
		}
	}
	return false
}

func (l *List[T]) Last() T {
	return l.Data[l.Len()-1]
}

func (l *List[T]) Len() int {
	return len(l.Data)
}
