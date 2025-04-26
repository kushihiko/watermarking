package bitset

import (
	"errors"
	"strings"
)

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

func NewBitSetFromString(bitsString string) (*BitSet, error) {
	bits := NewBitSet(len(bitsString))

	for i, b := range bitsString {
		switch b {
		case '1':
			bits.Set(i)
		case '0':
			continue
		default:
			return nil, errors.New("invalid character in bitset")
		}
	}

	return bits, nil
}

func (b *BitSet) Set(i int) {
	for i/8 >= len(b.bits) {
		b.bits = append(b.bits, 0)
	}

	if i >= b.Len() {
		b.length = i + 1
	}

	b.bits[i/8] |= 1 << (i % 8)
}

func (b *BitSet) Clear(i int) {
	b.bits[i/8] &^= 1 << (i % 8)
}

func (b *BitSet) Get(i int) bool {
	for i/8 >= len(b.bits) {
		b.bits = append(b.bits, 0)
	}

	if i >= b.Len() {
		b.length = i + 1
	}

	return (b.bits[i/8] & (1 << (i % 8))) != 0
}

func (b *BitSet) Len() int {
	return b.length
}

func (b *BitSet) String() string {
	var result strings.Builder

	for i := range b.Len() {
		if b.Get(i) {
			result.WriteRune('1')
		} else {
			result.WriteRune('0')
		}
	}

	return result.String()
}

func (b *BitSet) Add(x bool) {
	if x {
		b.Set(b.Len())
	} else {
		b.Get(b.Len())
	}
}
