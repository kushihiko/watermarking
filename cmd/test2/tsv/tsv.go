package tsv

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type TesseractUnit struct {
	Level      int
	Page       int
	Block      int
	Paragraph  int
	Line       int
	Word       int
	Left       int
	Top        int
	Width      int
	Height     int
	Confidence float64
	Text       string
}

func ParseTesseractTSV(path string) ([]TesseractUnit, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	var rows []TesseractUnit
	_, err = reader.Read() // Пропускаем заголовок
	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 12 {
			continue
		}

		level, _ := strconv.Atoi(record[0])
		page, _ := strconv.Atoi(record[1])
		block, _ := strconv.Atoi(record[2])
		paragraph, _ := strconv.Atoi(record[3])
		line, _ := strconv.Atoi(record[4])
		word, _ := strconv.Atoi(record[5])
		left, _ := strconv.Atoi(record[6])
		top, _ := strconv.Atoi(record[7])
		width, _ := strconv.Atoi(record[8])
		height, _ := strconv.Atoi(record[9])
		confidence, _ := strconv.ParseFloat(record[10], 64)
		text := record[11]

		rows = append(rows, TesseractUnit{
			Level:      level,
			Page:       page,
			Block:      block,
			Paragraph:  paragraph,
			Line:       line,
			Word:       word,
			Left:       left,
			Top:        top,
			Width:      width,
			Height:     height,
			Confidence: confidence,
			Text:       text,
		})
	}

	return rows, nil
}
