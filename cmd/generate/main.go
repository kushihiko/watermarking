package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/signintech/gopdf"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Загружаем слова из файла
	words, err := loadWords("russian_words_utf8.txt")
	if err != nil {
		panic(fmt.Sprintf("Не удалось загрузить слова: %v", err))
	}

	outputDir := "generated_pdfs"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Не удалось создать папку: %v", err))
	}

	for i := 0; i < 1000; i++ {
		text := generateRandomText(words, 5+rand.Intn(10)) // от 5 до 15 абзацев
		err := createPDF(fmt.Sprintf("%s/doc_%d.pdf", outputDir, i), text)
		if err != nil {
			fmt.Printf("Ошибка при создании документа %d: %v\n", i, err)
		}
	}
}

func loadWords(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	return words, scanner.Err()
}

func generateRandomText(words []string, paragraphs int) string {
	var builder strings.Builder
	capitalizeNext := true
	for i := 0; i < paragraphs; i++ {
		wordCount := 30 + rand.Intn(50) // от 30 до 80 слов на абзац
		for j := 0; j < wordCount; j++ {
			word := words[rand.Intn(len(words))]
			if capitalizeNext {
				word = strings.Title(word)
				capitalizeNext = false
			}
			builder.WriteString(word)
			r := rand.Float64()
			if r < 0.8 {
				builder.WriteString(" ")
			} else if r < 0.9 {
				builder.WriteString(", ")
			} else {
				builder.WriteString(". ")
				capitalizeNext = true
			}
			// 10% шанс завершить абзац досрочно
			if rand.Float64() < 0.05 {
				builder.WriteString("\n\n")
				capitalizeNext = true
			}
		}
		builder.WriteString("\n\n")
		capitalizeNext = true
	}
	return builder.String()
}

func createPDF(filename, text string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("Arial", "Arial.ttf")
	if err != nil {
		return fmt.Errorf("ошибка добавления шрифта: %w", err)
	}
	err = pdf.SetFont("Arial", "", 14)
	if err != nil {
		return fmt.Errorf("ошибка установки шрифта: %w", err)
	}

	margin := 40.0
	pageWidth := 595.0 // A4 ширина в gopdf
	usableWidth := pageWidth - 2*margin
	x := margin
	y := margin
	lineHeight := 20.0

	for _, paragraph := range strings.Split(text, "\n") {
		words := strings.Fields(paragraph)
		line := ""
		for _, word := range words {
			testLine := line
			if testLine != "" {
				testLine += " "
			}
			testLine += word
			w, _ := pdf.MeasureTextWidth(testLine)
			if w > usableWidth {
				pdf.SetX(x)
				pdf.SetY(y)
				pdf.Cell(nil, line)
				y += lineHeight
				line = word
			} else {
				line = testLine
			}
			if y+margin > gopdf.PageSizeA4.H {
				break
			}
		}
		if line != "" {
			pdf.SetX(x)
			pdf.SetY(y)
			pdf.Cell(nil, line)
			y += lineHeight
		}
		y += lineHeight // дополнительный отступ между абзацами
		if y+margin > gopdf.PageSizeA4.H {
			break
		}
	}

	return pdf.WritePdf(filename)
}
