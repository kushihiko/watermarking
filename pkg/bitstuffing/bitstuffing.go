package bitstuffing

import (
	"errors"
	"strings"
	"watermarking/pkg/bitset"
)

type BitStuffing struct {
	markerLength int
}

func NewBitStuffing(markerLength int) (*BitStuffing, error) {
	if markerLength < 3 {
		return nil, errors.New("marker length must be at least 3")
	}

	return &BitStuffing{
		markerLength: markerLength,
	}, nil

}

func (b *BitStuffing) Encode(set *bitset.BitSet) (*bitset.BitSet, error) {
	marker := "0" + strings.Repeat("1", b.markerLength-2) + "0"

	result, err := bitset.NewBitSetFromString(marker)
	if err != nil {
		return nil, err
	}

	onesInARaw := 0
	for i := range set.Len() {
		if set.Get(i) {
			onesInARaw++
		} else {
			onesInARaw = 0
		}

		result.Add(set.Get(i))

		if onesInARaw == b.markerLength-3 {
			result.Add(false)
			onesInARaw = 0
		}
	}

	return result, nil
}

func (b *BitStuffing) Decode(bset *bitset.BitSet) (*bitset.BitSet, error) {
	markerString := "0" + strings.Repeat("1", b.markerLength-2) + "0"
	marker, err := bitset.NewBitSetFromString(markerString)
	if err != nil {
		return nil, err
	}

	isMarker := true
	for i := 0; i < b.markerLength; i++ {
		if bset.Get(i) != marker.Get(i) {
			isMarker = false
			break
		}
	}

	iterator := 0
	if isMarker {
		iterator += b.markerLength
	}

	result := bitset.NewBitSet(0)
	onesInARaw := 0
	for ; iterator < bset.Len(); iterator++ {
		if bset.Get(iterator) {
			onesInARaw++
		} else {
			onesInARaw = 0
		}

		result.Add(bset.Get(iterator))

		if onesInARaw == b.markerLength-3 {
			iterator++
			onesInARaw = 0
			continue
		}
	}

	return result, nil
}
