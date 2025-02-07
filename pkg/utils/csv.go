package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func CSV2Excel(csvFile, excelFile string) error {

	cFile, err := os.Open(csvFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer cFile.Close()

	csvFilePath, _ := filepath.Abs(csvFile)
	dirPath := filepath.Dir(csvFilePath)
	cFileName := filepath.Base(csvFile)
	ext := filepath.Ext(csvFile)
	// Set the sheet name
	sheetName := strings.TrimSuffix(cFileName, ext)
	xlFile := filepath.Join(dirPath, excelFile)

	fmt.Println("csv file path: ", csvFilePath)

	// Read the CSV data
	reader := csv.NewReader(cFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("excel file path: ", xlFile)

	// Create a new Excel file
	// f := excelize.NewFile()
	f, err := excelize.OpenFile(xlFile)
	if err != nil {
		fmt.Println("File do not exits: Creating a new file")
		f = excelize.NewFile()
	} else {
		f.DeleteSheet(sheetName)
	}

	index, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println("Cannot create a sheet", err)
		return err
	}

	// Write the data to the sheet
	for i, record := range records {
		for j, value := range record {
			// Convert numeric values to float64
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				fmt.Println("Error", err)
			}
			if num, err := strconv.ParseFloat(value, 64); err == nil {
				f.SetCellValue(sheetName, cellName, num)
			} else {
				// Set the value and the style for non-numeric values
				f.SetCellValue(sheetName, cellName, value)
			}
		}
	}

	f.SetActiveSheet(index)
	// Save the file
	if err := f.SaveAs(xlFile); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Converted excel file stored at: ", xlFile)
	return nil
}

func CSV2Json(cFile string) (string, error) {
	csvFile, err := os.Open(cFile)
	if err != nil {
		fmt.Println("Open Error", err)
		return "", err
	}
	defer csvFile.Close()

	// get the file name and exetention for csvFile
	csvFilePath, _ := filepath.Abs(cFile)
	dirPath := filepath.Dir(csvFilePath)
	cFileName := filepath.Base(cFile)
	ext := filepath.Ext(cFile)
	fileName := strings.TrimSuffix(cFileName, ext)
	jFile := filepath.Join(dirPath, fileName)

	// Create a new JSON file
	jsonFile, err := os.Create(jFile + ".json")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer jsonFile.Close()

	// Read the CSV data
	reader := csv.NewReader(csvFile)
	reader.Comma = ','            // Set the delimiter to a comma
	headers, err := reader.Read() // Read the headers
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Loop through the records and convert them to a JSON object
	records := []map[string]string{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		// Convert the record to a map[string]string
		recordMap := make(map[string]string)
		for i, value := range record {
			if value != "" {
				// fmt.Println("Calue: ", value)
				recordMap[headers[i]] = value
			}

		}
		records = append(records, recordMap)
	}

	// Write the JSON object to the file
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println("Converted JSON file stored at:", jFile+".json")
	// fmt.Println("Converted JSON file stored at: ", string(jsonData))
	return jFile + ".json", nil
}
