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

	useCase, err := usecase.NewUseCase(cfg)
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
	case "extract":
		_, err := useCase.Extract()
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Неизвестная команда. Используйте embed или extract.")
		os.Exit(1)
	}
}
