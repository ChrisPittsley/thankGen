package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		errOut(fmt.Errorf("usage: thankGen template"))
	}
	t, err := os.ReadFile(os.Args[1])
	if err != nil {
		errOut(err)
	}
	template := string(t)
	records, err := csv.NewReader(os.Stdin).ReadAll()
	if err != nil {
		errOut(err)
	}
	for _, record := range records {
		fill := make([]interface{}, len(record))
		for i := range record {
			fill[i] = record[i]
		}
		fmt.Printf(template+"\n", fill...)
	}
}

func errOut(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	os.Exit(1)
}
