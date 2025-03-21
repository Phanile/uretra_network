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

	return b.Verify(*b.Signature)
}
