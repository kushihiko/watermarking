package main

import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	inputPDF := "input.pdf"
	outputPattern := "page-%d.png"

	err := convertPDFToImages(inputPDF, outputPattern)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("PDF успешно конвертирован в изображения")
	}
}

func convertPDFToImages(inputPath, outputPattern string) error {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Установим плотность DPI перед чтением
	mw.SetResolution(300, 300)

	// Загружаем все страницы PDF
	if err := mw.ReadImage(inputPath); err != nil {
		return fmt.Errorf("не удалось прочитать PDF: %w", err)
	}

	// Конвертация каждой страницы в PNG
	numPages := mw.GetNumberImages()
	for i := 0; i < int(numPages); i++ {
		mw.SetIteratorIndex(i)

		// К преобразованному изображению применим flatten (если есть прозрачность)
		page := mw.GetImage()
		page.SetImageFormat("png")

		outPath := fmt.Sprintf(outputPattern, i)
		if err := page.WriteImage(outPath); err != nil {
			return fmt.Errorf("не удалось сохранить %s: %w", outPath, err)
		}
		page.Destroy()
		fmt.Println("Сохранено:", outPath)
	}

	return nil
}
