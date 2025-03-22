package core

type Validator interface {
	ValidateBlock(*Block) bool
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(blockchain *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: blockchain,
	}
}

func (bv *BlockValidator) ValidateBlock(b *Block) bool {
	if bv.bc.HasBlock(b.Header.Height) {
		return false
	}

	if b.Header.Height != bv.bc.Height()+1 {
		return false
	}

	prevHeader, err := bv.bc.GetHeader(b.Header.Height - 1)
	hash := BlockHasher{}.Hash(prevHeader)

	if hash != b.Header.PrevBlockHash {
		return false
	}

	if err != nil {
		return false
	}

	return b.Verify()
}
