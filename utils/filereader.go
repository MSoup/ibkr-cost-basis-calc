package utils

import (
	"encoding/csv"
	"fmt"
	"io"
)

// ReadCSV reads a CSV file and returns a map where the key is the first column

type GroupedCSVRecords map[string][][]string

func ReadCSV(filepath io.Reader) GroupedCSVRecords {
	fmt.Printf("> Reading CSV\n")
	reader := csv.NewReader(filepath)

	csvMap := make(GroupedCSVRecords)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("> Loaded trades into memory")
			break
		}
		mapKey := record[0]
		// Make csvMap entry if header doesn't exist
		_, ok := csvMap[mapKey]
		if !ok {
			fmt.Println("> Creating new key:", mapKey)
			csvMap[mapKey] = make([][]string, 0)
		}
		csvMap[mapKey] = append(csvMap[mapKey], record)
	}

	return csvMap
}
