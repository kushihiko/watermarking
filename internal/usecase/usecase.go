package usecase

import (
	"fmt"
	"image"
	"strconv"
	"watermarking/pkg/bitset"
	"watermarking/pkg/painter"
	"watermarking/pkg/parser"
)

type UseCase struct {
	conv                                                              Converter
	prs                                                               Parser
	shft                                                              Shifter
	bstuff                                                            BitStuffing
	destroy                                                           func() error
	embedPDFSrc, embedPDFDst, extractPDFDst, outputPattern, tmpFolder string
	printBoxes                                                        bool
}

type Parser interface {
	Image(imagePath string) ([]parser.BoundingBox, error)
}

type Shifter interface {
	Normalize(boxes []image.Rectangle)
	Encrypt(boxes []image.Rectangle, bits bitset.BitSet)
	Decrypt(boxes []image.Rectangle) (bitset.BitSet, []float64)
}

type BitStuffing interface {
	PrepareMark(mark uint32) (*bitset.BitSet, error)
	Encode(set *bitset.BitSet) (*bitset.BitSet, error)
	RecoverRecords(detected []bitset.BitSet) *bitset.BitSet
	Decode(bset *bitset.BitSet) (*bitset.BitSet, error)
}

type Converter interface {
	PDFToImage(pdfPath string) ([]image.Image, error)
	ImagesToPDF(imagePaths []string, outputPath string) error
	Destroy() error
}

func NewUseCase(shft Shifter, bstuff BitStuffing, prs Parser, conv Converter, embedPDFSrc, embedPDFDst, extractPDFDst, outputPattern, tmpFolder string, printBoxes bool) (*UseCase, error) {
	return &UseCase{
		conv: conv,
		prs:  prs,
		shft: shft,
		destroy: func() error {
			return conv.Destroy()
		},
		bstuff:        bstuff,
		embedPDFSrc:   embedPDFSrc,
		embedPDFDst:   embedPDFDst,
		extractPDFDst: extractPDFDst,
		outputPattern: outputPattern,
		tmpFolder:     tmpFolder,
		printBoxes:    printBoxes,
	}, nil
}

func (uc *UseCase) Embed(mark uint32) error {
	imgs, err := uc.conv.PDFToImage(uc.embedPDFSrc)
	if err != nil {
		return err
	}

	encodedMark, err := uc.bstuff.PrepareMark(mark)
	if err != nil {
		return err
	}

	newImagePaths := make([]string, len(imgs))
	for i := range len(imgs) {
		imagePath := fmt.Sprintf(uc.tmpFolder+"/"+uc.outputPattern, i)
		newImagePath := fmt.Sprintf(uc.tmpFolder+"/newImages/"+uc.outputPattern, i)

		if err = painter.SaveImage(imgs[i], imagePath); err != nil {
			return err
		}

		wordsBoxes, err := uc.prs.Image(imagePath)
		if err != nil {
			return err
		}

		//if err = painter.DeleteImage(imagePath); err != nil {
		//	return err
		//}

		newBoxes := make([]image.Rectangle, len(wordsBoxes))
		oldBoxes := make([]image.Rectangle, len(wordsBoxes))
		words := make([]string, len(wordsBoxes))
		for j, wordsBox := range wordsBoxes {
			newBoxes[j] = wordsBox.Box
			oldBoxes[j] = wordsBox.Box
			words[j] = wordsBox.Text
		}

		uc.shft.Normalize(newBoxes)
		uc.shft.Encrypt(newBoxes, *encodedMark)

		pnt := painter.NewPainter(imgs[i].Bounds().Dx(), imgs[i].Bounds().Dy(), uc.printBoxes)
		newImg, err := pnt.Rearrange(imgs[i], oldBoxes, newBoxes, words)
		// newImg, err := pnt.DrawBoxes(imgs[i], oldBoxes)
		if err != nil {
			return err
		}

		err = painter.SaveImage(newImg, newImagePath)
		if err != nil {
			return err
		}
		newImagePaths[i] = newImagePath
	}

	err = uc.conv.ImagesToPDF(newImagePaths, uc.extractPDFDst)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) Extract() (marker uint32, err error) {
	imgs, err := uc.conv.PDFToImage(uc.extractPDFDst)
	if err != nil {
		return 0, err
	}

	dectipted := make([]bitset.BitSet, len(imgs))
	for i := range len(imgs) {
		imagePath := fmt.Sprintf(uc.tmpFolder+"/"+uc.outputPattern, i)
		newImagePath := fmt.Sprintf(uc.tmpFolder+"/newImages/"+uc.outputPattern, i)

		if err = painter.SaveImage(imgs[i], imagePath); err != nil {
			return 0, err
		}

		wordsBoxes, err := uc.prs.Image(imagePath)
		if err != nil {
			return 0, err
		}

		//if err = painter.DeleteImage(imagePath); err != nil {
		//	return 0, err
		//}

		boxes := make([]image.Rectangle, len(wordsBoxes))
		for j := range wordsBoxes {
			boxes[j] = wordsBoxes[j].Box
		}

		// debug
		pnt := painter.NewPainter(imgs[i].Bounds().Dx(), imgs[i].Bounds().Dy(), uc.printBoxes)
		newImg, err := pnt.DrawBoxes(imgs[i], boxes)
		if err != nil {
			return 0, err
		}
		err = painter.SaveImage(newImg, newImagePath)
		if err != nil {
			return 0, err
		}

		bSetDecr, _ := uc.shft.Decrypt(boxes)
		dectipted[i] = bSetDecr
	}

	recoveredMark := uc.bstuff.RecoverRecords(dectipted)
	mark, err := uc.bstuff.Decode(recoveredMark)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseUint(mark.String(), 2, 32) // база 2 (двоичная), 32 бита
	if err != nil {
		// return 0, err
		return 0, nil
	}

	return uint32(value), nil
}

func (uc *UseCase) Destroy() error {
	return uc.destroy()
}
