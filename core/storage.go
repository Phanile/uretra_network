package core

type Storage interface {
	Put(*Block) error
	Get(height uint32) error
}

type MemoryStorage struct {
	blocks []*Block
}

func (ms *MemoryStorage) Put(b *Block) error {
	ms.blocks = append(ms.blocks, b)
	return nil
}

func (ms *MemoryStorage) Get(h uint32) error {
	return nil
}
