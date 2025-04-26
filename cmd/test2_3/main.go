package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"image"
	"watermarking/pkg/bitset"
	"watermarking/pkg/bitstuffing"
	"watermarking/pkg/converter"
	"watermarking/pkg/painter"
	"watermarking/pkg/parser"
	"watermarking/pkg/shifter"
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
		"rus+eng",
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

	bitstf, err := bitstuffing.NewBitStuffing(5)
	if err != nil {
		panic(err)
	}

	bset, err := bitset.NewBitSetFromString("11111011011111101010111110")
	if err != nil {
		panic(err)
	}

	newBSet, err := bitstf.Encode(bset)
	if err != nil {
		panic(err)
	}

	fmt.Println("HAHAHA:     ", bset.String())
	fmt.Println("HAHAHA:", newBSet.String())

	decodedBSET, err := bitstf.Decode(newBSet)
	if err != nil {
		panic(err)
	}

	fmt.Println("HAHAHA:     ", decodedBSET.String())
}
