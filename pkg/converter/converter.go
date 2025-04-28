package converter

import (
	"bytes"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
)

type Converter struct {
	mw                                    *imagick.MagickWand
	gaussianBlurRadius, gaussianBlurSigma float64
	noise                                 imagick.NoiseType
	noiceOffset                           float64
	destroy                               func()
}

func NewConverter(gaussianBlurRadius, gaussianBlurSigma float64, noise imagick.NoiseType, noiceOffset float64) *Converter {
	imagick.Initialize()
	mw := imagick.NewMagickWand()

	return &Converter{
		mw:                 mw,
		gaussianBlurRadius: gaussianBlurRadius,
		gaussianBlurSigma:  gaussianBlurSigma,
		noise:              noise,
		noiceOffset:        noiceOffset,
		destroy: func() {
			mw.Destroy()
			imagick.Terminate()
		},
	}
}

func (c *Converter) Clear() {
	c.mw.Clear()
}

func (c *Converter) PDFToImage(pdfPath string) ([]image.Image, error) {
	if err := c.mw.SetResolution(200, 200); err != nil {
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
		defer page.Destroy()

		width := page.GetImageWidth()
		height := page.GetImageHeight()

		if width > 2500 || height > 3500 {
			var coef float64
			if width > 2500 {
				coef = 2500.0 / float64(width)
			} else {
				coef = 3500.0 / float64(height)
			}
			newW := uint(float64(width) * coef)
			newH := uint(float64(height) * coef)

			if err := page.ResizeImage(newW, newH, imagick.FILTER_LANCZOS); err != nil {
				return imgs, err
			}
		}

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

	err := finalWand.GaussianBlurImage(c.gaussianBlurRadius, c.gaussianBlurSigma)
	if err != nil {
		return err
	}
	err = finalWand.AddNoiseImage(c.noise, c.noiceOffset)
	if err != nil {
		return err
	}

	if err := finalWand.SetFormat("pdf"); err != nil {
		return err
	}

	if err := finalWand.WriteImages(outputPath, true); err != nil {
		return err
	}

	return nil
}

func (c *Converter) Destroy() error {
	c.destroy()
	return nil
}
