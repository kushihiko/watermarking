package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	PDFPath       string `yaml:"pdf_path"`
	PDFName       string `yaml:"pdf_name"`
	TmpFolder     string `yaml:"tmp_folder"`
	OutputFolder  string `yaml:"output_folder"`
	PrintBoxes    bool   `yaml:"print_boxes"`
	OutputPattern string `yaml:"output_pattern"`
	FontPath      string `yaml:"font_path"`
	Language      string `yaml:"language"`
	WhiteList     string `yaml:"whitelist"`
	BlackList     string `yaml:"blacklist"`
	Watermark     string `yaml:"watermark"`
	Shift         int    `yaml:"shift"`
	MarkerLength  int    `yaml:"marker_length"`
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
