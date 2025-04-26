package converter

import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
)

const TmpFolder = "tmp/"

type Converter struct {
	mw      *imagick.MagickWand
	Destroy func()
}

func NewConverter() *Converter {
	imagick.Initialize()
	mw := imagick.NewMagickWand()

	return &Converter{
		mw: mw,
		Destroy: func() {
			imagick.Terminate()
			mw.Destroy()
		},
	}
}

func (c *Converter) SetGaussianBlur(radius, sigma float64) error {
	return c.mw.GaussianBlurImage(radius, sigma)
}

func (c *Converter) AddNoise(noise imagick.NoiseType, offset float64) error {
	return c.mw.AddNoiseImage(noise, offset)
}

func (c *Converter) convertPDFToImage(pdfPath string) error {
	if err := c.mw.SetResolution(300, 300); err != nil {
		return err
	}

	if err := c.mw.ReadImage(pdfPath); err != nil {
		return fmt.Errorf("не удалось прочитать PDF: %w", err)
	}

	outputPattern := "page-%d.png"
	numPages := c.mw.GetNumberImages()
	for i := 0; i < int(numPages); i++ {
		c.mw.SetIteratorIndex(i)

		page := c.mw.GetImage()
		if err := page.SetImageFormat("png"); err != nil {
			return err
		}

		outPath := fmt.Sprintf(TmpFolder+outputPattern, i)
		if err := page.WriteImage(outPath); err != nil {
			return fmt.Errorf("не удалось сохранить %s: %w", outPath, err)
		}

		page.Destroy()
	}

	return nil
}
