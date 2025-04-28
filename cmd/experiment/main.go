package main

import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"log"
	"os"
	"watermarking/internal/config"
	"watermarking/internal/usecase"
	"watermarking/pkg/bitstuffing"
	"watermarking/pkg/converter"
	gocvparser "watermarking/pkg/parser/gocv"
	"watermarking/pkg/shifter"
)

type Experiment struct {
	Name        string
	BlurRadius  float64
	BlurSigma   float64
	NoiseType   imagick.NoiseType
	NoiseOffset float64
}

func main() {
	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// Создаём зависимости
	shft := shifter.NewShifter(cfg.Embed.Shift)
	bstuff, err := bitstuffing.NewBitStuffing(cfg.Embed.MarkerLength)
	if err != nil {
		log.Fatal(err)
	}
	prs := gocvparser.NewParser(10)

	baseMark := uint32(123456) // Тестовая метка

	// Описываем сценарии экспериментов
	experiments := []Experiment{
		{"Без искажений", 0, 0, 0, 0},
		//{"Блюр слабый", 1.0, 0.5, 0, 0},
		//{"blur", 2.0, 1.0, imagick.NOISE_UNDEFINED, 0},
		//{"noise", 0, 0, imagick.NOISE_GAUSSIAN, 0.5},
		//{"high_noise", 0, 0, imagick.NOISE_GAUSSIAN, 0.5},
		//{"blur_noise", 1.5, 0.7, imagick.NOISE_GAUSSIAN, 0.5},
	}

	fmt.Println("| Эксперимент | Среднее совпадение (%) |")
	fmt.Println("|-------------|------------------------|")

	for _, exp := range experiments {
		conv := converter.NewConverter(exp.BlurRadius, exp.BlurSigma, exp.NoiseType, exp.NoiseOffset)

		successes := 0
		total := 0

		dst := fmt.Sprintf("generated_pdfs_embedded/%s", exp.Name)
		os.MkdirAll(dst, os.ModePerm)

		for i := 0; i < 1; i++ {
			srcPath := fmt.Sprintf("../generate/generated_pdfs/doc_%d.pdf", i)
			dstPath := fmt.Sprintf("generated_pdfs_embedded/%s/doc_%d.pdf", exp.Name, i)
			extractPath := fmt.Sprintf("generated_pdfs_embedded/%s/doc_%d.pdf", exp.Name, i)

			useCase, err := usecase.NewUseCase(shft, bstuff, prs, conv, srcPath, dstPath, extractPath, cfg.OutputPattern, cfg.TmpFolder, cfg.Embed.PrintBoxes)
			if err != nil {
				log.Fatal(err)
			}

			// Встраиваем метку
			err = useCase.Embed(baseMark)
			if err != nil {
				log.Fatalf("Ошибка встраивания для документа %d: %v", i, err)
			}
			// useCase.Destroy()

			conv.Clear()
			// Извлекаем метку
			useCase, err = usecase.NewUseCase(shft, bstuff, prs, conv, dstPath, dstPath, extractPath, cfg.OutputPattern, cfg.TmpFolder, cfg.Embed.PrintBoxes)
			if err != nil {
				log.Fatal(err)
			}

			extractedMark, err := useCase.Extract()
			if err != nil {
				log.Fatalf("Ошибка извлечения для документа %d: %v", i, err)
			}
			// useCase.Destroy()

			if extractedMark == baseMark {
				successes++
			}
			total++

			conv.Clear()
			fmt.Println(exp.Name, "test:", i)
		}

		accuracy := float64(successes) / float64(total) * 100
		fmt.Printf("| %s | %.2f%% |\n", exp.Name, accuracy)
		conv.Destroy()
	}
}
