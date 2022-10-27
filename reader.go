package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"
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
	Template *template.Template
}

func New(tmpl string, cfg Config) (Magi, error) {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return Magi{}, err
	}
	if cfg.StrictMode {
		t.Option("missingkey=error")
	} else {
		t.Option("missingkey=zero")
	}

	magi := Magi{
		Config:   cfg,
		Template: t,
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

		err = magi.Template.Execute(w, data)
		if err != nil {
			return fmt.Errorf("line %d: error %w\n", line, err)
		}

		// add new line
		_, err = w.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

type header []string

func keyVariants(index int, key string) []string {
	// FIXME: this leads to converting Foo+Bar and Foo-Bar to the
	// same varibales Foo_Bar, Foo_Bar.
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, "+", "_")
	key = strings.ReplaceAll(key, "-", "_")
	key = strings.ReplaceAll(key, ".", "_")
	// FIXME: what about combination of underscore lower case and upper case
	// like key FoObAr? Maybe add proper lookup to template but currently we only
	// have option with function which makes templates very complex.
	return []string{
		key,
		strings.ToLower(key),
		strings.ToUpper(key),
		// support special column number like _1
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

	for i, column := range header {
		variants := keyVariants(i, column)
		fmt.Printf("{{.%s}} ", variants[0])
	}
	fmt.Println()
}
