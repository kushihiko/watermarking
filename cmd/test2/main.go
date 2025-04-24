package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/otiai10/gosseract/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
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

	//for _, wordsBox := range wordsBoxes {
	//	psr.Word(wordsBox)
	//}

	// Преобразование всех слов
	for i, wordsBox := range wordsBoxes {
		//letterBox := psr.Word(wordsBox)
		//for _, box := range letterBox {
		//	dc.SetColor(color.RGBA{255, 0, 0, 255})
		//	dc.DrawRectangle(float64(box.Box.Min.X), float64(box.Box.Min.Y), float64(box.Box.Dx()), float64(box.Box.Dy()))
		//	dc.Stroke()
		//}
		//fmt.Println(wordsBox)

		dc.SetColor(color.RGBA{0, 255, 0, 255})
		dc.DrawRectangle(float64(wordsBox.Box.Min.X), float64(wordsBox.Box.Min.Y), float64(wordsBox.Box.Dx()), float64(wordsBox.Box.Dy()))
		dc.Stroke()

		if i != 0 {
			fmt.Println("DIFF:", wordsBoxes[i-1].Box.Max.X-wordsBoxes[i].Box.Min.X)
		}
	}

	// Сохраняем изображение
	outFile, _ := os.Create(outPath)
	defer outFile.Close()
	png.Encode(outFile, dc.Image())

	fmt.Println("Результат сохранён в", outPath)
}
