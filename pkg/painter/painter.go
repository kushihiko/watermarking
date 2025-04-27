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
	printBoxes    bool
}

func NewPainter(width, height int, printBoxes bool) *Painter {
	return &Painter{
		width:      width,
		height:     height,
		printBoxes: printBoxes,
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
	for i := range boxes {
		wordImg := image.NewRGBA(image.Rect(0, 0, boxes[i].Dx(), boxes[i].Dy()))
		draw.Draw(wordImg, wordImg.Bounds(), oldImage, oldBoxes[i].Min, draw.Src)

		newDC.DrawImage(wordImg, boxes[i].Min.X, boxes[i].Min.Y)

		if p.printBoxes {
			newDC.SetRGBA(1, 0, 0, 0.8)
			newDC.SetLineWidth(2)
			newDC.DrawRectangle(float64(boxes[i].Min.X), float64(boxes[i].Min.Y), float64(boxes[i].Dx()), float64(boxes[i].Dy()))
			newDC.Stroke()
		}
	}

	return newDC.Image(), nil
}

func (p *Painter) DrawBoxes(img image.Image, boxes []image.Rectangle) (image.Image, error) {
	newDC := gg.NewContextForImage(img)
	newDC.SetRGB(1, 1, 1)
	//newDC.Clear()

	for i := range boxes {
		newDC.SetRGBA(1, 0, 0, 0.8)
		newDC.SetLineWidth(2)
		newDC.DrawRectangle(float64(boxes[i].Min.X), float64(boxes[i].Min.Y), float64(boxes[i].Dx()), float64(boxes[i].Dy()))
		newDC.Stroke()

		boxImg := image.NewRGBA(image.Rect(0, 0, boxes[i].Dx(), boxes[i].Dy()))
		draw.Draw(boxImg, boxImg.Bounds(), img, boxes[i].Min, draw.Src)
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
