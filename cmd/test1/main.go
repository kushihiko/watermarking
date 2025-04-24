package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/golang/freetype"
	"golang.org/x/net/html"
)

func main() {
	// Читаем HOCR
	file, err := os.Open("test/output.hocr")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		panic(err)
	}

	// Создаём белое изображение
	imgWidth, imgHeight := 2000, 1000
	dst := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Загружаем TTF-шрифт
	fontBytes, err := os.ReadFile("Arial.ttf")
	if err != nil {
		panic(err)
	}
	ttf, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetFontSize(18)
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.NewUniform(color.Black))

	// Рекурсивный парсинг слов
	var processNode func(*html.Node)
	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			class := getAttr(n, "class")
			if class == "ocrx_word" {
				text := getNodeText(n)
				bbox := parseBbox(getAttr(n, "title"))
				if bbox != nil && text != "" {
					pt := freetype.Pt(bbox.Min.X, bbox.Max.Y-4) // -4: подгон по базовой линии
					_, err := c.DrawString(text, pt)
					if err != nil {
						fmt.Println("draw error:", err)
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			processNode(child)
		}
	}
	processNode(doc)

	// Сохраняем изображение
	outFile, _ := os.Create("output_reconstructed.png")
	defer outFile.Close()
	png.Encode(outFile, dst)
	fmt.Println("✅ Готово: output_reconstructed.png")
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func getNodeText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		text := getNodeText(child)
		if text != "" {
			return text
		}
	}
	return ""
}

func parseBbox(title string) *image.Rectangle {
	parts := strings.Split(title, ";")
	for _, part := range parts {
		if strings.HasPrefix(strings.TrimSpace(part), "bbox") {
			fields := strings.Fields(part)
			if len(fields) == 5 {
				x1, _ := strconv.Atoi(fields[1])
				y1, _ := strconv.Atoi(fields[2])
				x2, _ := strconv.Atoi(fields[3])
				y2, _ := strconv.Atoi(fields[4])
				rect := image.Rect(x1, y1, x2, y2)
				return &rect
			}
		}
	}
	return nil
}
