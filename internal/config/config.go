package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TmpFolder     string  `yaml:"tmp_folder"`
	OutputFolder  string  `yaml:"output_folder"`
	OutputPattern string  `yaml:"output_pattern"`
	Embed         Embed   `yaml:"embed"`
	Extract       Extract `yaml:"extract"`
}

type Embed struct {
	PDFSrc       string `yaml:"pdf_src"`
	PDFDst       string `yaml:"pdf_dst"`
	PrintBoxes   bool   `yaml:"print_boxes"`
	Shift        int    `yaml:"shift"`
	MarkerLength int    `yaml:"marker_length"`
}

type Extract struct {
	PDFSrc string `yaml:"pdf_src"`
}

func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
