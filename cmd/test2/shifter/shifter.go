package shifter

import (
	"fmt"
	"image"
	"watermarking/cmd/test2/bitset"
)

type Shifter struct{}

func (sh *Shifter) Normalize(words []image.Rectangle) {
	var totalGap int
	var gapCount int
	for i := 1; i < len(words); i++ {
		gap := words[i].Min.X - words[i-1].Max.X
		if gap > 0 {
			totalGap += gap
			gapCount++
		}
	}

	avgGap := 0
	if gapCount > 0 {
		avgGap = totalGap / gapCount
	}

	cursorX := words[0].Min.X
	prevX := words[0].Min.X
	for i := range words {
		if prevX > words[i].Min.X {
			cursorX = words[i].Min.X
		}

		dx := words[i].Dx()
		words[i].Min.X = cursorX
		words[i].Max.X = cursorX + dx

		cursorX += dx + avgGap
		prevX = words[i].Min.X
	}
}

func (sh *Shifter) Encrypt(boxes []image.Rectangle, shift int, bits bitset.BitSet) {
	cursorX := boxes[0].Min.X
	prevX := boxes[0].Min.X

	bitNumber := 0

	for i := range boxes {
		if prevX > boxes[i].Min.X {
			cursorX = boxes[i].Min.X
			bitNumber--
		}

		if bitNumber >= bits.Len() {
			bitNumber = 0
		}

		dx := boxes[i].Dx()
		boxes[i].Min.X = cursorX
		boxes[i].Max.X = cursorX + dx

		if bits.Get(bitNumber) {
			cursorX += shift
		}

		cursorX += dx + shift
		prevX = boxes[i].Min.X
		bitNumber++
	}
}

// Generated
func (sh *Shifter) Decrypt(boxes []image.Rectangle) (bitset.BitSet, []float64) {
	if len(boxes) < 2 {
		return *bitset.NewBitSet(0), nil
	}

	var gaps []int
	mp := make(map[int]int)
	for i := 1; i < len(boxes); i++ {
		gap := boxes[i].Min.X - boxes[i-1].Max.X
		if gap > 0 {
			mp[gap]++
			gaps = append(gaps, gap)
		}
	}

	fmt.Println("GAPS:", mp)

	if len(gaps) == 0 {
		return *bitset.NewBitSet(0), nil
	}

	// Find min and max gaps
	minGap, maxGap := gaps[0], gaps[0]
	for _, g := range gaps {
		if g < minGap {
			minGap = g
		}
		if g > maxGap {
			maxGap = g
		}
	}

	// Assume two clusters: small (0) and large (1)
	center0 := minGap
	center1 := maxGap

	bits := bitset.NewBitSet(len(gaps))
	confidences := make([]float64, len(gaps))

	for i, g := range gaps {
		diff0 := abs(g - center0)
		diff1 := abs(g - center1)

		if diff1 < diff0 {
			bits.Set(i)
		}

		totalDiff := diff0 + diff1
		if totalDiff != 0 {
			confidences[i] = 1.0 - (float64(min(diff0, diff1)) / float64(totalDiff))
		} else {
			confidences[i] = 1.0
		}
	}

	return *bits, confidences
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
