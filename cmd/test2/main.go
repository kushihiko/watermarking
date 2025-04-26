package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/otiai10/gosseract/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/draw"
	"image/png"
	"os"
	"watermarking/cmd/test2/parser"
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
	// tsvPath := "/Users/kushihiko/Projects/watermarking/test/output.tsv"
	imagePath := "/Users/kushihiko/Projects/watermarking/test/2.png"
	outPath := "output_boxes.png"

	fontSize := 14.0

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

	//// Читаем TSV
	//tesseractUnits, err := tsv.ParseTesseractTSV(tsvPath)
	//if err != nil {
	//	panic(err)
	//}

	psr, err := parser.NewParser(
		"rus",
		gosseract.RIL_WORD,
		"АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя",
		"-.,:;",
		fontPath,
	)
	if err != nil {
		panic(err)
	}
	defer psr.Close()

	wordsBoxes, err := psr.Image(imagePath)
	if err != nil {
		panic(err)
	}

	// Новый холст для выравненного текста
	//newDC := gg.NewContext(width, height)
	//newDC.SetRGB(1, 1, 1)
	//newDC.Clear()
	//
	//// Считаем средний пробел между словами
	//var totalGap int
	//var gapCount int
	//for i := 1; i < len(wordsBoxes); i++ {
	//	prev := wordsBoxes[i-1].Box
	//	curr := wordsBoxes[i].Box
	//	gap := curr.Min.X - prev.Max.X
	//	if gap > 0 && wordsBoxes[i-1].LineNum == wordsBoxes[i].LineNum {
	//		totalGap += gap
	//		gapCount++
	//	}
	//}
	//
	//avgGap := 0
	//if gapCount > 0 {
	//	avgGap = totalGap / gapCount
	//}
	//
	//// Начальная позиция
	//cursorX := wordsBoxes[0].Box.Min.X
	//prevX := wordsBoxes[0].Box.Min.X
	//for _, word := range wordsBoxes {
	//	if prevX > word.Box.Min.X {
	//		cursorX = word.Box.Min.X
	//	}
	//
	//	wordImg := image.NewRGBA(image.Rect(0, 0, word.Box.Dx(), word.Box.Dy()))
	//	draw.Draw(wordImg, wordImg.Bounds(), img, word.Box.Min, draw.Src)
	//
	//	// Копируем слово на новое место
	//	newDC.DrawImage(wordImg, cursorX, word.Box.Min.Y)
	//
	//	// Переход к следующему слову
	//	cursorX += word.Box.Dx() + avgGap
	//	prevX = word.Box.Min.X
	//}
	//
	//// Сохраняем новое изображение
	//outAlignedFile, _ := os.Create("output_aligned.png")
	//defer outAlignedFile.Close()
	//png.Encode(outAlignedFile, newDC.Image())
	//
	//// Сохраняем изображение
	//outFile, _ := os.Create(outPath)
	//defer outFile.Close()
	//png.Encode(outFile, dc.Image())
	//
	//fmt.Println("Результат сохранён в", outPath)
}
