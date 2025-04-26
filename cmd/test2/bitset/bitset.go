package bitset

type BitSet struct {
	bits   []byte
	length int
}

func NewBitSet(size int) *BitSet {
	return &BitSet{
		bits:   make([]byte, (size+7)/8),
		length: size,
	}
}

func (b *BitSet) Set(i int) {
	for i/8 >= len(b.bits) {
		b.bits = append(b.bits, 0)
		b.length = i + 1
	}

	b.bits[i/8] |= (1 << (i % 8))
}

func (b *BitSet) Clear(i int) {
	b.bits[i/8] &^= (1 << (i % 8))
}

func (b *BitSet) Get(i int) bool {
	if i/8 >= len(b.bits) {
		return false
	}

	return (b.bits[i/8] & (1 << (i % 8))) != 0
}

func (b *BitSet) Len() int {
	return b.length
}
