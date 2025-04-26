package converter

import (
	"bytes"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
)

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
			mw.Destroy()
			imagick.Terminate()
		},
	}
}

func (c *Converter) SetGaussianBlur(radius, sigma float64) error {
	return c.mw.GaussianBlurImage(radius, sigma)
}

func (c *Converter) AddNoise(noise imagick.NoiseType, offset float64) error {
	return c.mw.AddNoiseImage(noise, offset)
}

func (c *Converter) PDFToImage(pdfPath string) ([]image.Image, error) {
	if err := c.mw.SetResolution(300, 300); err != nil {
		return nil, err
	}

	if err := c.mw.ReadImage(pdfPath); err != nil {
		return nil, fmt.Errorf("не удалось прочитать PDF: %w", err)
	}

	var imgs []image.Image
	numPages := c.mw.GetNumberImages()
	for i := 0; i < int(numPages); i++ {
		c.mw.SetIteratorIndex(i)

		page := c.mw.GetImage()
		if err := page.SetImageFormat("png"); err != nil {
			return imgs, err
		}

		blob, err := page.GetImageBlob()
		if err != nil {
			return imgs, err
		}
		img, err := blobToImage(blob)
		if err != nil {
			return imgs, err
		}
		imgs = append(imgs, img)

		//outPath := fmt.Sprintf(c.imageFolder+"/"+c.outputPattern, i)
		//if err := page.WriteImage(outPath); err != nil {
		//	return uint(i), fmt.Errorf("не удалось сохранить %s: %w", outPath, err)
		//}

		page.Destroy()
	}

	return imgs, nil
}

func blobToImage(blob []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(blob))
	return img, err
}

func (c *Converter) ImagesToPDF(imagePaths []string, outputPath string) error {
	finalWand := imagick.NewMagickWand()
	defer finalWand.Destroy()

	for _, path := range imagePaths {
		tempWand := imagick.NewMagickWand()

		if err := tempWand.ReadImage(path); err != nil {
			return err
		}

		if err := finalWand.AddImage(tempWand); err != nil {
			return err
		}

		tempWand.Destroy()
	}

	if err := finalWand.SetFormat("pdf"); err != nil {
		return err
	}

	if err := finalWand.WriteImages(outputPath, true); err != nil {
		return err
	}

	return nil
}
