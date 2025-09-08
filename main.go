package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type config struct {
	out io.Writer
	template string
	records [][]string
}

func (cfg *config) setOutput(path string) error {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	cfg.out, err = os.Create(path)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *config) setTemplate(path string) error {
	var err error
	if cfg.template != "" {
		return fmt.Errorf("template already set, could not use '%s' as template", path)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("no valid template found in '%s'", path)
	}
	cfg.template = string(data)
	return nil
}

func (cfg *config) addTable(path string) error {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	table, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return err
	}
	cfg.records = append(cfg.records, table...)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		errOut(fmt.Errorf("not enough arguments"))
	}
	cfg, err := setup(os.Args[1:])
	if err != nil {
		errOut(err)
	}
	for _, record := range cfg.records {
		fill := make([]any, len(record))
		for i := range record {
			fill[i] = record[i]
		}
		fmt.Fprintf(cfg.out, cfg.template+"\n\n", fill...)
	}
}

func errOut(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	fmt.Fprint(os.Stderr, "usage: thankGen template.txt table.csv [-o output.txt]\n")
	os.Exit(1)
}

func setup(args []string) (config, error) {
	var cfg config
	var o bool
	cfg.out = os.Stdout
	for i := 0; i < len(args); i++ {
		if args[i] == "-o" {
			i++
			if i >= len(args) {
				return cfg, fmt.Errorf("-o specified without filename")
			}
			if o {
				return cfg, fmt.Errorf("-o specified twice, output already set")
			}
			if err := cfg.setOutput(args[i]); err != nil {
				return cfg, err
			}
			o = true
			continue
		}
		switch filepath.Ext(args[i]) {
		case ".csv":
			if err := cfg.addTable(args[i]); err != nil {
				return cfg, err
			}
		default:
			if err := cfg.setTemplate(args[i]); err != nil {
				return cfg, err
			}
		}
	}
	if cfg.template == "" {
		return cfg, fmt.Errorf("no template specified")
	}
	if len(cfg.records) == 0 {
		return cfg, fmt.Errorf("no table specified")
	}
	return cfg, nil
}
