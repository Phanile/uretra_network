package core

type Blockchain struct {
	store     Storage
	headers   []*Header
	validator Validator
}

func NewBlockchain(genesis *Block) *Blockchain {
	bc := &Blockchain{
		headers: []*Header{},
		store:   &MemoryStorage{},
	}

	bc.validator = NewBlockValidator(bc)

	err := bc.addBlockWithoutValidation(genesis)

	if err != nil {
		panic(err)
	}

	return bc
}

func (bc *Blockchain) addBlock(b *Block) bool {
	if bc.validator.ValidateBlock(b) {
		err := bc.addBlockWithoutValidation(b)

		if err != nil {
			return false
		}

		return true
	}

	return false
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)

	return bc.store.Put(b)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) SetValidator(val Validator) {
	bc.validator = val
}

func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}
