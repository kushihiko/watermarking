package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
	"watermarking/pkg/tsv"
)

type WordBox struct {
	Text      string
	Left, Top int
	Width     int
	Height    int
}

type LetterBox struct {
	Char       rune
	XMin, YMin int
	XMax, YMax int
}

func main() {
	// Пути
	fontPath := "/System/Library/Fonts/Supplemental/Times New Roman.ttf"
	tsvPath := "/Users/kushihiko/Projects/watermarking/test/output.tsv"
	imagePath := "image-0.png"
	outPath := "output_boxes.png"

	fontSize := 24.0

	// Загружаем шрифт
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     300,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	defer face.Close()

	// Загружаем изображение
	imgFile, _ := os.Open(imagePath)
	img, _, _ := image.Decode(imgFile)
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	imgFile.Close()

	// Создаём холст
	dc := gg.NewContext(width, height)
	dc.DrawImage(img, 0, 0)

	// Читаем TSV
	tesseractUnits, err := tsv.ParseTesseractTSV(tsvPath)
	if err != nil {
		panic(err)
	}

	// Преобразование всех слов
	for _, unit := range tesseractUnits {
		boxes := getLetterBoxes(word, face)
		for _, box := range boxes {
			dc.SetColor(color.RGBA{255, 0, 0, 255})
			dc.DrawRectangle(float64(box.XMin), float64(box.YMin), float64(box.XMax-box.XMin), float64(box.YMax-box.YMin))
			dc.Stroke()
		}
	}

	// Сохраняем изображение
	outFile, _ := os.Create(outPath)
	defer outFile.Close()
	png.Encode(outFile, dc.Image())

	fmt.Println("Результат сохранён в", outPath)
}

func getLetterBoxes(word WordBox, face font.Face) []LetterBox {
	boxes := []LetterBox{}
	text := word.Text
	if len(text) == 0 {
		return boxes
	}

	charRunes := []rune(text)
	advances := make([]float64, len(charRunes))
	var totalAdvance float64

	for i, r := range charRunes {
		advance, ok := face.GlyphAdvance(r)
		if !ok {
			advance = 0
		}
		advances[i] = float64(advance.Round())
		totalAdvance += advances[i]
	}

	if totalAdvance == 0 {
		return boxes
	}

	scale := float64(word.Width) / totalAdvance
	x := float64(word.Left)

	for i, r := range charRunes {
		w := advances[i] * scale

		bounds, _, ok := face.GlyphBounds(r)
		if !ok {
			continue
		}

		// Преобразуем границы символа в пиксели и масштабируем под word.Height
		glyphHeight := float64(bounds.Max.Y - bounds.Min.Y)
		totalFontHeight := float64(face.Metrics().Height)
		scaleY := float64(word.Height) / totalFontHeight

		yMin := float64(word.Top) + float64(face.Metrics().Ascent-bounds.Max.Y)*scaleY
		yMax := yMin + glyphHeight*scaleY

		box := LetterBox{
			Char: r,
			XMin: int(x),
			YMin: int(yMin),
			XMax: int(x + w),
			YMax: int(yMax),
		}
		boxes = append(boxes, box)
		x += w
	}

	return boxes
}
