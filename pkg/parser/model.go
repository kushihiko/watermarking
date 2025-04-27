package parser

import "image"

type BoundingBox struct {
	Box        image.Rectangle
	Confidence float64
	Text       string
}
