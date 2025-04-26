package usecase

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"image"
	"watermarking/internal/config"
	"watermarking/pkg/bitset"
	"watermarking/pkg/bitstuffing"
	"watermarking/pkg/converter"
	"watermarking/pkg/painter"
	"watermarking/pkg/parser"
	"watermarking/pkg/shifter"
)

type UseCase struct {
	cfg     *config.Config
	conv    *converter.Converter
	prs     *parser.Parser
	shft    *shifter.Shifter
	destroy func() error
}

func NewUseCase(cfg *config.Config) UseCase {
	conv := converter.NewConverter()
	return UseCase{
		cfg:  cfg,
		conv: conv,
		prs: parser.NewParser(
			cfg.Language,
			gosseract.RIL_WORD,
		),
		shft: shifter.NewShifter(),
		destroy: func() error {
			conv.Destroy()
			return nil
		},
	}
}

func (uc *UseCase) Embed(mark uint32) error {
	imgs, err := uc.conv.PDFToImage(uc.cfg.PDFPath)
	if err != nil {
		return err
	}

	encodedMark, err := uc.prepareMark(mark)
	if err != nil {
		return err
	}

	newImagePaths := make([]string, len(imgs))
	for i := range len(imgs) {
		imagePath := fmt.Sprintf(uc.cfg.TmpFolder+"/"+uc.cfg.OutputPattern, i)
		newImagePath := fmt.Sprintf(uc.cfg.TmpFolder+"/newImages/"+uc.cfg.OutputPattern, i)

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
			words[j] = wordsBox.Word
		}

		// TODO: fix normilaze: floating letters
		uc.shft.Normalize(newBoxes)
		uc.shft.Encrypt(newBoxes, uc.cfg.Shift, *encodedMark)

		pnt := painter.NewPainter(imgs[i].Bounds().Dx(), imgs[i].Bounds().Dy())
		newImg, err := pnt.Rearrange(imgs[i], oldBoxes, newBoxes, words)
		if err != nil {
			return err
		}
		err = painter.SaveImage(newImg, newImagePath)
		if err != nil {
			return err
		}
		newImagePaths[i] = newImagePath
	}

	err = uc.conv.ImagesToPDF(newImagePaths, uc.cfg.OutputFolder+"/"+uc.cfg.PDFName)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) prepareMark(mark uint32) (*bitset.BitSet, error) {
	bitstf, err := bitstuffing.NewBitStuffing(uc.cfg.MarkerLength)
	if err != nil {
		return nil, err
	}

	bset, err := bitset.NewBitSetFromString(fmt.Sprintf("%b", mark))
	if err != nil {
		return nil, err
	}

	encodedBSet, err := bitstf.Encode(bset)
	if err != nil {
		return nil, err
	}

	return encodedBSet, nil
}

func (uc *UseCase) Extract() (mark uint32, err error) {
	return 10, nil
}

func (uc *UseCase) Destroy() error {
	return uc.destroy()
}
