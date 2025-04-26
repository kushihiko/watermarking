package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"watermarking/internal/config"
	"watermarking/internal/usecase"
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
		fmt.Println("Использование: watermark_cli [embed|extract] --config=config.yaml --storage=0 --event=0")
		os.Exit(1)
	}
	cmd := flag.Arg(0)

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// TODO: init
	mark := uint32(eventID)

	useCase := usecase.NewUseCase(cfg)
	defer useCase.Destroy()

	switch cmd {
	case "embed":
		err = useCase.Embed(mark)
	case "extract":
		fmt.Println("Извлечение водяной метки...")
		// TODO: реализовать extract
	default:
		fmt.Println("Неизвестная команда. Используйте embed или extract.")
		os.Exit(1)
	}
}
