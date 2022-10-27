package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

var (
	ErrHeaderExpected = errors.New("header expected")
	// TODO: add line information
	ErrNoSuchColumn = errors.New("no such column")
)

type Config struct {
	// In strict mode errors:
	// 1. Not enough columns
	// 2. No column in template
	StrictMode bool
}

type Magi struct {
	Config   Config
	Template string
}

func New(tmpl string, cfg Config) (Magi, error) {
	magi := Magi{
		Config:   cfg,
		Template: tmpl,
	}
	return magi, nil
}

func (magi *Magi) ReadAndExecute(r io.Reader, w io.Writer) error {
	csvReader := csv.NewReader(r)
	if !magi.Config.StrictMode {
		csvReader.FieldsPerRecord = -1
	}
	header, err := readHeader(csvReader)
	if err != nil {
		return err
	}

	for line := 1; ; line++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("line %d: error %w\n", line, err)
		}

		data, err := header.mapWithColumns(magi.Config.StrictMode, record)
		if err != nil {
			return fmt.Errorf("line %d: error %w\n", line, err)
		}

		result := Execute(magi.Template, data)

		_, err = w.Write([]byte(result + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

type header []string

func keyVariants(index int, key string) []string {
	return []string{
		key,
		fmt.Sprintf("_%d", index+1),
	}
}

func (h header) mapWithColumns(
	strictMode bool, values []string,
) (map[string]string, error) {
	m := make(map[string]string)
	for i := range h {
		var value string
		if i < len(values) {
			value = values[i]
		} else if strictMode {
			return nil, ErrNoSuchColumn
		}
		for _, key := range keyVariants(i, h[i]) {
			m[key] = value
		}
	}
	return m, nil
}

func readHeader(r *csv.Reader) (header, error) {
	record, err := r.Read()
	if err == io.EOF {
		return nil, ErrHeaderExpected
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

func PrintExample(r io.Reader) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1
	header, err := readHeader(csvReader)
	if err != nil {
		panic(err)
	}

	fmt.Print("   '")
	for _, column := range header {
		fmt.Printf("{%s} ", column)
	}
	fmt.Println("'")
}
