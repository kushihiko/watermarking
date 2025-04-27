package gocvparser

import (
	"errors"
	"gocv.io/x/gocv"
	"image"
	"sort"
	"watermarking/pkg/parser"
)

type Parser struct {
	minArea int
}

func NewParser(minArea int) *Parser {
	return &Parser{
		minArea: minArea,
	}
}

func (p *Parser) Close() {}

func (p *Parser) Image(imagePath string) ([]parser.BoundingBox, error) {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return nil, errors.New("gocv: image read error")
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	bin := gocv.NewMat()
	defer bin.Close()

	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	gocv.Threshold(gray, &bin, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

	contours := gocv.FindContours(bin, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	var chars []image.Rectangle
	var avgDx, avgDy int
	for i := 0; i < contours.Size(); i++ {
		rect := gocv.BoundingRect(contours.At(i))
		area := rect.Dx() * rect.Dy()
		if area < p.minArea {
			continue
		}
		avgDx += rect.Bounds().Dx()
		avgDy += rect.Bounds().Dy()
		chars = append(chars, rect)
	}
	avgDx /= len(chars)
	avgDy /= len(chars)

	maxXGap := int(float64(avgDx) / 2.3)
	maxYGap := avgDy * 2

	sort.Slice(chars, func(i, j int) bool {
		if abs(chars[i].Min.Y-chars[j].Min.Y) > maxYGap {
			return chars[i].Min.Y < chars[j].Min.Y
		}
		return chars[i].Min.X < chars[j].Min.X
	})

	//debugImg := img.Clone()
	//defer debugImg.Close()
	//
	//for i, r := range chars {
	//	gocv.Rectangle(&debugImg, r, color.RGBA{0, 255, 0, 0}, 2)
	//
	//	pt := image.Pt(r.Min.X, r.Min.Y-5)
	//	gocv.PutText(&debugImg, fmt.Sprintf("%d", i), pt, gocv.FontHersheySimplex, 0.3, color.RGBA{255, 0, 0, 0}, 1)
	//}
	//
	//gocv.IMWrite("debug_output.png", debugImg)

	var boxes []parser.BoundingBox
	var currentBox *image.Rectangle

	for _, rect := range chars {
		if currentBox == nil {
			box := rect
			currentBox = &box
			continue
		}

		// Проверяем, близки ли прямоугольники по X и Y
		xGap := rect.Min.X - currentBox.Max.X
		yGap := abs(rect.Min.Y - currentBox.Min.Y)

		if xGap <= maxXGap && yGap <= maxYGap {
			// Расширяем текущий прямоугольник
			currentBox.Max.X = rect.Max.X
			currentBox.Max.Y = max(currentBox.Max.Y, rect.Max.Y)
			currentBox.Min.Y = min(currentBox.Min.Y, rect.Min.Y)
		} else {
			// Добавляем законченную группу
			boxes = append(boxes, parser.BoundingBox{
				Box:        *currentBox,
				Confidence: 1.0,
				Text:       "",
			})
			box := rect
			currentBox = &box
		}
	}

	// Добавляем последнюю группу
	if currentBox != nil {
		boxes = append(boxes, parser.BoundingBox{
			Box:        *currentBox,
			Confidence: 1.0,
			Text:       "",
		})
	}

	return boxes, nil
}

// Вспомогательные функции
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
