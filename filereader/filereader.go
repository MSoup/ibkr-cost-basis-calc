package filereader

import (
	"encoding/csv"
	"fmt"
	"io"
)

func ReadCSV(filepath io.Reader) map[string][][]string {
	fmt.Printf("> Reading CSV\n")
	reader := csv.NewReader(filepath)

	csvMap := make(map[string][][]string)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("> Loaded CSV into memory")
			break
		}

		// Make csvMap entry if header doesn't exist
		_, ok := csvMap[record[0]]
		if !ok {
			csvMap[record[0]] = make([][]string, 0)
		}
		csvMap[record[0]] = append(csvMap[record[0]], record)
	}

	return csvMap
}
