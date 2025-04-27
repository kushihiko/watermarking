package tesseract

import (
	"github.com/otiai10/gosseract/v2"
	"watermarking/pkg/parser"
)

type Parser struct {
	languages string
	level     gosseract.PageIteratorLevel
	whiteList string
	blackList string
}

func NewParser(languages string, level gosseract.PageIteratorLevel) *Parser {
	//fontBytes, err := os.ReadFile(fontPath)
	//if err != nil {
	//	return nil, err
	//}
	//ft, err := opentype.Parse(fontBytes)
	//if err != nil {
	//	return nil, err
	//}
	//face, err := opentype.NewFace(ft, &opentype.FaceOptions{
	//	Size:    9,
	//	DPI:     300,
	//	Hinting: font.HintingFull,
	//})
	//if err != nil {
	//	return nil, err
	//}

	return &Parser{
		languages: languages,
		level:     level,
	}
}

func (p *Parser) Image(imagePath string) ([]parser.BoundingBox, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetLanguage(p.languages)
	if err != nil {
		return nil, err
	}

	//err = client.SetWhitelist(p.whiteList)
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = client.SetBlacklist(p.blackList)
	//if err != nil {
	//	return nil, err
	//}

	err = client.SetImage(imagePath)
	if err != nil {
		return nil, err
	}

	boxes, err := client.GetBoundingBoxes(p.level)
	if err != nil {
		return nil, err
	}

	result := make([]parser.BoundingBox, len(boxes))
	for i := range boxes {
		result[i] = parser.BoundingBox{
			Box:        boxes[i].Box,
			Confidence: boxes[i].Confidence,
			Text:       boxes[i].Word,
		}
	}

	return result, nil
}

//type letter struct {
//	height         fixed.Int26_6
//	width          fixed.Int26_6
//	baselineOffset fixed.Int26_6
//	char           rune
//}
//
//type LetterBoundingBox struct {
//	Box                                image.Rectangle
//	Letter                             rune
//	Confidence                         float64
//	BlockNum, ParNum, LineNum, WordNum int
//}

//func (p *Parser) Word(wordsBox gosseract.BoundingBox) []LetterBoundingBox {
//	var minBaselineOffset fixed.Int26_6
//	var sumFontWidth fixed.Int26_6
//	var maxHeight fixed.Int26_6
//	var letters []letter
//
//	var prevChar rune
//	for i, char := range wordsBox.Word {
//		var height, width, baselineOffset fixed.Int26_6
//		bounds, advance, ok := p.face.GlyphBounds(char)
//		// p.face.Glyph()
//		if !ok {
//			return nil
//		}
//
//		width = advance
//		height = bounds.Max.Y - bounds.Min.Y
//		baselineOffset = bounds.Min.Y
//
//		var kern fixed.Int26_6
//		if i != 0 {
//			kern = p.face.Kern(prevChar, char)
//		} else {
//			minBaselineOffset = baselineOffset
//		}
//
//		fmt.Println("KERN", Int26_6ToFloat64(kern))
//
//		sumFontWidth += advance + kern
//
//		if baselineOffset < minBaselineOffset {
//			minBaselineOffset = baselineOffset
//		}
//
//		letters = append(letters, letter{
//			height:         height,
//			width:          width,
//			baselineOffset: baselineOffset,
//			char:           char,
//		})
//		prevChar = char
//	}
//
//	stretchFactor := float64(wordsBox.Box.Dx()) / Int26_6ToFloat64(sumFontWidth)
//	// stretchFactor := float64(wordsBox.Box.Dy()) / Int26_6ToFloat64(maxHeight)
//	fmt.Println("MAXHeight", Int26_6ToFloat64(maxHeight))
//
//	currentMin := image.Point{
//		X: wordsBox.Box.Min.X,
//		Y: wordsBox.Box.Min.Y,
//	}
//
//	fmt.Println("WORD", wordsBox.Box)
//	// fmt.Println("STRETCH", stretchFactor)
//
//	lettersBox := make([]LetterBoundingBox, len(letters))
//	for i, letter := range letters {
//		box := image.Rectangle{
//			Min: image.Point{
//				X: currentMin.X,
//				Y: currentMin.Y + int(math.Round(Int26_6ToFloat64(letter.baselineOffset-minBaselineOffset)*stretchFactor)),
//			},
//			Max: image.Point{
//				X: currentMin.X + int(math.Round(Int26_6ToFloat64(letter.width)*stretchFactor)),
//				Y: currentMin.Y + int(math.Round(Int26_6ToFloat64(letter.baselineOffset-minBaselineOffset+letter.height)*stretchFactor)),
//			},
//		}
//
//		var kern fixed.Int26_6
//		if i != 0 {
//			kern = p.face.Kern(letters[i-1].char, letters[i].char)
//		}
//
//		currentMin.X += int(math.Round(Int26_6ToFloat64(letter.width+kern) * stretchFactor))
//		lettersBox[i] = LetterBoundingBox{
//			Box:        box,
//			Letter:     letter.char,
//			Confidence: wordsBox.Confidence,
//			BlockNum:   wordsBox.BlockNum,
//			ParNum:     wordsBox.ParNum,
//			LineNum:    wordsBox.LineNum,
//			WordNum:    wordsBox.WordNum,
//		}
//
//		fmt.Println(lettersBox[i].Box)
//	}
//
//	return lettersBox
//}
//
//func Int26_6ToFloat64(x fixed.Int26_6) float64 {
//	return float64(x) / 64.0
//}
