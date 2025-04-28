package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"watermarking/internal/config"
	"watermarking/internal/usecase"
	"watermarking/pkg/bitstuffing"
	"watermarking/pkg/converter"
	gocvparser "watermarking/pkg/parser/gocv"
	"watermarking/pkg/shifter"
)

type UseCase interface {
	Embed(mark uint32) error
	Extract() (mark uint32, err error)
	Destroy() error
}

func main() {
	var configPath string
	var storageID int
	var eventID int

	flag.StringVar(&configPath, "config", "config.yaml", "Путь к YAML конфигу")
	flag.IntVar(&storageID, "storage", 0, "Идентификатор места хранения событий")
	flag.IntVar(&eventID, "event", 0, "Идентификатор события")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Использование: watermark_cli --config=config.yaml --storage=0 --event=0 [embed|extract]")
		os.Exit(1)
	}
	cmd := flag.Arg(0)

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// TODO: init
	mark := uint32(eventID)

	shft := shifter.NewShifter(cfg.Embed.Shift)
	bstuff, err := bitstuffing.NewBitStuffing(cfg.Embed.MarkerLength)
	if err != nil {
		log.Fatal(err)
	}
	prs := gocvparser.NewParser(10)
	conv := converter.NewConverter(0, 0, 0, 0)

	useCase, err := usecase.NewUseCase(shft, bstuff, prs, conv, cfg.Embed.PDFSrc, cfg.Embed.PDFDst, cfg.Extract.PDFSrc, cfg.OutputPattern, cfg.TmpFolder, cfg.Embed.PrintBoxes)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = useCase.Destroy()
		if err != nil {
			log.Println(err)
		}
	}()

	switch cmd {
	case "embed":
		err = useCase.Embed(mark)
		if err != nil {
			log.Fatal(err)
		}
	case "extract":
		mark, err := useCase.Extract()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(mark)
	default:
		fmt.Println("Неизвестная команда. Используйте embed или extract.")
		os.Exit(1)
	}
}
