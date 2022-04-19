package file

import (
	"encoding/csv"
	"io"
	"log"
)

func ReadPaging(file io.Reader, offset, limit int) ([][]string, error) {
	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	var line int
	records := make([][]string, 0)
	for {
		rec, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line++
		if line <= offset {
			continue
		}
		records = append(records, rec)
		if line >= limit {
			break
		}
	}
	return records, nil
}

func ReadAll(reader io.Reader, skipHeader bool) ([][]string, error) {
	// read csv values using csv.Reader
	csvReader := csv.NewReader(reader)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if skipHeader {
		return data[1:], nil
	}
	return data, nil
}

func WriteAll(file io.Writer, records [][]string) error {
	csvWriter := csv.NewWriter(file)
	err := csvWriter.WriteAll(records)
	if err != nil {
		return err
	}
	return nil
}

func Write(file io.Writer, record []string) error {
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	err := csvWriter.Write(record)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
