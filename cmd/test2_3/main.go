package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"image"
	"watermarking/cmd/test2/bitset"
	"watermarking/cmd/test2/converter"
	"watermarking/cmd/test2/painter"
	"watermarking/cmd/test2/parser"
	"watermarking/cmd/test2/shifter"
)

func main() {
	pdfPath := "test.pdf"
	imageFolder := "tmp/"
	outputFolder := "output/"
	outputPattern := "page-%d.png"
	fontPath := "/System/Library/Fonts/Supplemental/Times New Roman.ttf"

	conv := converter.NewConverter()
	defer conv.Destroy()

	imgs, err := conv.PDFToImage(pdfPath)
	if err != nil {
		panic(err)
	}

	prs, err := parser.NewParser(
		"rus",
		gosseract.RIL_WORD,
		"АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя",
		"-.,:;",
		fontPath,
	)
	if err != nil {
		panic(err)
	}

	var shft shifter.Shifter
	for i := range len(imgs) {
		imagePath := fmt.Sprintf(imageFolder+"/"+outputPattern, i)

		if err = painter.SaveImage(imgs[i], imagePath); err != nil {
			panic(err)
		}

		wordsBoxes, err := prs.Image(imagePath)
		if err != nil {
			panic(err)
		}

		//if err = painter.DeleteImage(imagePath); err != nil {
		//	panic(err)
		//}

		newBoxes := make([]image.Rectangle, len(wordsBoxes))
		oldBoxes := make([]image.Rectangle, len(wordsBoxes))
		words := make([]string, len(wordsBoxes))
		for j, wordsBox := range wordsBoxes {
			newBoxes[j] = wordsBox.Box
			oldBoxes[j] = wordsBox.Box
			words[j] = wordsBox.Word
		}

		shft.Normalize(newBoxes)
		btset, err := bitset.NewBitSetFromString("01110101101100")
		if err != nil {
			panic(err)
		}
		shft.Encrypt(newBoxes, 4, *btset)

		decryptBitSet, _ := shft.Decrypt(newBoxes)
		fmt.Println(decryptBitSet.String())
		//fmt.Println("CONF:", confidences)

		pnt := painter.NewPainter(imgs[i].Bounds().Dx(), imgs[i].Bounds().Dy())
		newImg, err := pnt.Rearrange(imgs[i], oldBoxes, newBoxes, words)
		if err != nil {
			panic(err)
		}

		outputPath := fmt.Sprintf(outputFolder+"/"+outputPattern, i)
		if err = painter.SaveImage(newImg, outputPath); err != nil {
			panic(err)
		}
	}
}
