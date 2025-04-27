package bitstuffing

import (
	"errors"
	"fmt"
	"strings"
	"watermarking/pkg/bitset"
)

type BitStuffing struct {
	markerLength int
	marker       *bitset.BitSet
}

func NewBitStuffing(markerLength int) (*BitStuffing, error) {
	if markerLength < 3 {
		return nil, errors.New("marker length must be at least 3")
	}

	markerString := "0" + strings.Repeat("1", markerLength-2) + "0"
	marker, _ := bitset.NewBitSetFromString(markerString)

	return &BitStuffing{
		markerLength: markerLength,
		marker:       marker,
	}, nil

}

func (b *BitStuffing) PrepareMark(mark uint32) (*bitset.BitSet, error) {
	bset, err := bitset.NewBitSetFromString(fmt.Sprintf("%b", mark))
	if err != nil {
		return nil, err
	}

	encodedBSet, err := b.Encode(bset)
	if err != nil {
		return nil, err
	}

	return encodedBSet, nil
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

func (b *BitStuffing) RecoverRecords(detected []bitset.BitSet) *bitset.BitSet {
	var allRecords []string
	marker := "0" + strings.Repeat("1", b.markerLength-2) + "0"

	for i := range detected {
		bitStr := detected[i].String()
		indexes := findAllMarkers(bitStr, marker)

		for j := 0; j < len(indexes); j++ {
			start := indexes[j]
			end := len(bitStr)
			if j+1 < len(indexes) {
				end = indexes[j+1]
			}
			if start >= end {
				continue
			}

			allRecords = append(allRecords, bitStr[start:end])
		}
	}

	if len(allRecords) == 0 {
		return bitset.NewBitSet(0)
	}

	distribution := make(map[int]int)
	for i := range allRecords {
		distribution[len(allRecords[i])]++
	}

	mostCommonLen := 0
	maxCount := 0
	for length, count := range distribution {
		if count > maxCount {
			mostCommonLen = length
			maxCount = count
		}
	}

	var filteredRecords []string
	for i := range allRecords {
		if abs(len(allRecords[i])-mostCommonLen) <= b.markerLength {
			filteredRecords = append(filteredRecords, allRecords[i])
		}
	}

	votes := make([]int, mostCommonLen)
	counts := make([]int, mostCommonLen)

	for i := range filteredRecords {
		for j, char := range filteredRecords[i] {
			if j < mostCommonLen {
				if char == '1' {
					votes[j]++
				}
				counts[j]++
			}
		}
	}

	result := bitset.NewBitSet(mostCommonLen)
	for i := 0; i < mostCommonLen; i++ {
		if counts[i] > 0 && votes[i]*2 >= counts[i] {
			result.Set(i)
		}
	}

	return result
}

func findAllMarkers(s, marker string) []int {
	var indexes []int
	idx := strings.Index(s, marker)
	for idx != -1 {
		indexes = append(indexes, idx)
		idx = strings.Index(s[idx+1:], marker)
		if idx != -1 {
			idx += indexes[len(indexes)-1] + 1
		}
	}
	return indexes
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
