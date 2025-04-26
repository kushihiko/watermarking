package painter

import (
	"errors"
	"github.com/fogleman/gg"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type Painter struct {
	width, height int
}

func NewPainter(width, height int) *Painter {
	return &Painter{
		width:  width,
		height: height,
	}

}

func (p *Painter) Rearrange(oldImage image.Image, oldBoxes []image.Rectangle, boxes []image.Rectangle, words []string) (image.Image, error) {
	if len(words) != len(boxes) {
		return nil, errors.New("words and boxes are not equal")
	}

	newDC := gg.NewContext(p.width, p.height)
	newDC.SetRGB(1, 1, 1)
	newDC.Clear()

	// Начальная позиция
	for i, box := range boxes {
		wordImg := image.NewRGBA(image.Rect(0, 0, box.Dx(), box.Dy()))
		draw.Draw(wordImg, wordImg.Bounds(), oldImage, oldBoxes[i].Min, draw.Src)

		newDC.DrawImage(wordImg, box.Min.X, box.Min.Y)
	}

	return newDC.Image(), nil
}

func SaveImage(img image.Image, imagePath string) error {
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

func DeleteImage(imagePath string) error {
	return os.Remove(imagePath)
}
