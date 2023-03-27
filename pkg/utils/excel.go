package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"encoding/json"

	"github.com/xuri/excelize/v2"
)

// ExcelSheet2CSV function for one sheet
func ExcelSheet2CSV(excelFile, sheetName string) error {

	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
		}
	}()

	excelize.NewFile()
	xlFilePath, _ := filepath.Abs(excelFile)
	dirPath := filepath.Dir(xlFilePath)

	fmt.Println("Excel File Path: ", dirPath)

	csvFile := filepath.Join(dirPath, sheetName+".csv")

	cFile, err := os.Create(csvFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer cFile.Close()

	csvFilePath, _ := filepath.Abs(csvFile)

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	write2csv(rows, cFile)

	fmt.Println("CSV File Path: ", csvFilePath)
	return nil
}

// Excel2CSV function for one sheet
func Excel2CSV(excelFile string) error {

	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
		}
	}()

	xlFilePath, _ := filepath.Abs(excelFile)
	dirPath := filepath.Dir(xlFilePath)
	fmt.Println("Excel File Path: ", xlFilePath)

	sheets := f.GetSheetList()

	for i := 0; i < len(sheets); i++ {
		// fmt.Println(sheets[i])
		csvFile := filepath.Join(dirPath, sheets[i]+".csv")
		csvFilePath, _ := filepath.Abs(csvFile)

		cFile, err := os.Create(csvFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer cFile.Close()

		// fmt.Println("CSV File Path: ", csvFilePath)
		// Get all the rows in the Sheet1.
		rows, err := f.GetRows(sheets[i])
		if err != nil {
			fmt.Println(err)
			return err
		}

		write2csv(rows, cFile)

		fmt.Println("CSV File Path: ", csvFilePath)
	}
	return nil
}

func Excel2Json(excelFile string) (string, error) {
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		return "", err
	}

	xlFilePath, _ := filepath.Abs(excelFile)
	dirPath := filepath.Dir(xlFilePath)
	fmt.Println("Excel File Path: ", xlFilePath)

	// Get the name of the first sheet
	sheetName := f.GetSheetName(0)

	jsonFile := filepath.Join(dirPath, sheetName+".json")

	// Get the rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Create a slice to hold the JSON data
	data := []map[string]string{}
	var headers []string

	if len(rows) == 0 {
		// fmt.Println("The Excel sheet is empty")
		return "", errors.New("the excel sheet is empty")
	} else {
		// Get the headers from the first row
		fmt.Println("length: ", len(rows))
		headers = rows[0]
	}

	// Loop through the remaining rows
	for i := 1; i < len(rows); i++ {
		// Create a map to hold the row data
		rowData := make(map[string]string)

		// Loop through the cells in the row
		for j, cellValue := range rows[i] {
			// Add the cell value to the row data map with the header as the key
			rowData[headers[j]] = cellValue
		}

		// Add the row data map to the slice
		data = append(data, rowData)
	}

	// Convert the data slice to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Write the JSON data to a file
	err = WriteFile(jsonFile, jsonData)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println("Conversion complete")
	return jsonFile, nil
}

// WriteFile writes data to a file at the specified path
func WriteFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func write2csv(rows [][]string, cFile *os.File) {
	nbOfColums := 0
	for i, row := range rows {
		// fmt.Println("Row Length", len(row))
		for j, colCell := range row {

			if i < 1 {
				nbOfColums = j + 1
			}

			re := regexp.MustCompile("\\n")
			colCell = re.ReplaceAllString(colCell, " ")

			re = regexp.MustCompile("\\r")
			colCell = re.ReplaceAllString(colCell, " ")

			// fmt.Println("Lenght of Coloum Cell", len(colCell))
			if len(colCell) < 1000 {
				colCell = strings.Replace(colCell, ",", " ", -1)
				cFile.WriteString(colCell + ",")
			} else {
				fmt.Println("Including only 1000 charectors from  row and coloum", i, j)
				colCell = strings.Replace(colCell, ",", " ", -1)
				cFile.WriteString(colCell[0:1000] + ",")
			}
			// fmt.Println("Colons and colom", nbOfColums, len(row))

		}
		if nbOfColums > len(row) {
			cFile.WriteString(",")
		}
		// fmt.Println("After of Columns", j)

		cFile.WriteString("\n")
	}
}
