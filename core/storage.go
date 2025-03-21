package core

type Storage interface {
	Put(*Block) error
}

type MemoryStorage struct {
}

func (ms *MemoryStorage) Put(*Block) error {
	return nil
}
