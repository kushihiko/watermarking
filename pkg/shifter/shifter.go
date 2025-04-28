package shifter

import (
	"fmt"
	"image"
	"watermarking/pkg/bitset"
)

type Shifter struct {
	shift int
}

func NewShifter(shift int) *Shifter {
	return &Shifter{
		shift: shift,
	}
}

func (sh *Shifter) Normalize(boxes []image.Rectangle) {
	var totalGap int
	var gapCount int
	for i := 1; i < len(boxes); i++ {
		gap := boxes[i].Min.X - boxes[i-1].Max.X
		if gap > 0 {
			totalGap += gap
			gapCount++
		}
	}

	avgGap := 0
	if gapCount > 0 {
		avgGap = totalGap / gapCount
	}

	cursorX := boxes[0].Min.X
	prevMaxY := boxes[0].Max.Y
	for i := range boxes {
		if prevMaxY < boxes[i].Min.Y {
			cursorX = boxes[i].Min.X
		}

		dx := boxes[i].Dx()
		boxes[i].Min.X = cursorX
		boxes[i].Max.X = cursorX + dx

		cursorX += dx + avgGap
		prevMaxY = boxes[i].Max.Y
	}
}

func (sh *Shifter) Encrypt(boxes []image.Rectangle, bits bitset.BitSet) {
	cursorX := boxes[0].Min.X
	prevMaxY := boxes[0].Max.Y
	gap := boxes[1].Min.X - boxes[0].Max.X

	bitNumber := 0

	for i := range boxes {
		if prevMaxY < boxes[i].Min.Y {
			cursorX = boxes[i].Min.X
			bitNumber--
		}

		if bitNumber >= bits.Len() || bitNumber < 0 {
			bitNumber = 0
		}

		dx := boxes[i].Dx()
		boxes[i].Min.X = cursorX
		boxes[i].Max.X = cursorX + dx

		if bits.Get(bitNumber) {
			cursorX += sh.shift
		}

		cursorX += dx + sh.shift + gap
		prevMaxY = boxes[i].Max.Y
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

	fmt.Println(mp)

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
