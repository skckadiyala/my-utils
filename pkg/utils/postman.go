package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type postmanResults struct {
	Run struct {
		Meta struct {
			CollectionID   string `json:"collectionId"`
			CollectionName string `json:"collectionName"`
			Started        int64  `json:"started"`
			Completed      int64  `json:"completed"`
			Duration       int    `json:"duration"`
		} `json:"meta"`
		Summary struct {
			Iterations struct {
				Executed int `json:"executed"`
				Errors   int `json:"errors"`
			} `json:"iterations"`
			ExecutedRequests struct {
				Executed int `json:"executed"`
				Errors   int `json:"errors"`
			} `json:"executedRequests"`
			PrerequestScripts struct {
				Executed int `json:"executed"`
				Errors   int `json:"errors"`
			} `json:"prerequestScripts"`
			PostresponseScripts struct {
				Executed int `json:"executed"`
				Errors   int `json:"errors"`
			} `json:"postresponseScripts"`
			Tests struct {
				Executed int `json:"executed"`
				Failed   int `json:"failed"`
				Passed   int `json:"passed"`
				Skipped  int `json:"skipped"`
			} `json:"tests"`
			TimeStats struct {
				ResponseAverage           float64 `json:"responseAverage"`
				ResponseMin               int     `json:"responseMin"`
				ResponseMax               int     `json:"responseMax"`
				ResponseStandardDeviation float64 `json:"responseStandardDeviation"`
			} `json:"timeStats"`
		} `json:"summary"`
		Executions []struct {
			IterationCount  int `json:"iterationCount"`
			RequestExecuted struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				URL  struct {
					Protocol string   `json:"protocol"`
					Path     []string `json:"path"`
					Host     []string `json:"host"`
					Query    []any    `json:"query"`
					Variable []any    `json:"variable"`
				} `json:"url"`
				Method string `json:"method"`
			} `json:"requestExecuted"`
			Response struct {
				ID      string `json:"id"`
				Details struct {
					Name         string `json:"name"`
					Detail       string `json:"detail"`
					Code         int    `json:"code"`
					StandardName string `json:"standardName"`
				} `json:"_details"`
				Status          string `json:"status"`
				Code            int    `json:"code"`
				Cookies         []any  `json:"cookies"`
				ResponseTime    int    `json:"responseTime"`
				ResponseSize    int    `json:"responseSize"`
				DownloadedBytes int    `json:"downloadedBytes"`
			} `json:"response"`
			Tests []struct {
				Name  string `json:"name"`
				Error struct {
					Name    string `json:"name"`
					Index   int    `json:"index"`
					Test    string `json:"test"`
					Message string `json:"message"`
					Stack   string `json:"stack"`
				} `json:"error,omitempty"`
				Status string `json:"status"`
			} `json:"tests"`
			Errors []any `json:"errors"`
		} `json:"executions"`
		RunError any `json:"runError"`
	} `json:"run"`
}

func concatenateWithDelimiter(arr []string, delimiter string) string {
	return strings.Join(arr, delimiter)
}

func readPostmanResultsFromFile(filePath string) (postmanResults, error) {
	var results postmanResults
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return results, err
	}
	err = json.Unmarshal(file, &results)
	return results, err
}

func PostmanResults2Splunk(resultFile, splunkHost, splunkPort, userName, password, source, index string) error {

	results, err := readPostmanResultsFromFile(resultFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return err
	}

	splunkUrl := "https://" + splunkHost + ":" + splunkPort
	userData := []byte(userName + ":" + password)
	basicAuth := base64.StdEncoding.EncodeToString(userData)

	extractedData := make(map[string]string)
	extractedData["collectionName"] = results.Run.Meta.CollectionName
	extractedData["automationTool"] = "Postman"
	extractedData["jobName"] = os.Getenv("BUILD_DEFINITIONNAME")
	extractedData["runID"] = os.Getenv("BUILD_BUILDNUMBER")

	startDate := results.Run.Meta.Started

	//convert to started time to 2025-01-15 format

	startedTime := time.Unix(0, startDate*int64(time.Millisecond)).Format("2006-01-02")
	extractedData["testDate"] = startedTime
	extractedData["environment"] = os.Getenv("Environment")

	for num, execution := range results.Run.Executions {
		extractedData["requestName"] = execution.RequestExecuted.Name
		concatenateWithDelimiter(execution.RequestExecuted.URL.Host, "/")
		concatenateWithDelimiter(execution.RequestExecuted.URL.Path, "/")

		extractedData["url"] = execution.RequestExecuted.URL.Protocol + "://" + concatenateWithDelimiter(execution.RequestExecuted.URL.Host, ".") + "/" + concatenateWithDelimiter(execution.RequestExecuted.URL.Path, "/")
		extractedData["method"] = execution.RequestExecuted.Method
		extractedData["responseTime"] = fmt.Sprintf("%d", execution.Response.ResponseTime)
		extractedData["responseStatus"] = execution.Response.Status
		extractedData["responseCode"] = fmt.Sprintf("%d", execution.Response.Code)
		extractedData["responseDetails"] = execution.Response.Details.Detail

		for _, test := range execution.Tests {
			extractedData["testSatus"] = test.Status
			extractedData["testName"] = test.Name

			if test.Status == "failed" {
				extractedData["testError"] = test.Error.Message
			} else {
				extractedData["testError"] = ""
			}

			// fmt.Println(execution.RequestExecuted.Name + ":" + test.Name)

			data, err := json.MarshalIndent(extractedData, "", "  ")
			if err != nil {
				return err
			}

			fmt.Printf("Test Result %v : %v \n", num+1, string(data))
			PostResults(splunkUrl, source, index, basicAuth, data)

			// testNum := fmt.Sprintf("%d", noOfTest)
			// nameext := fmt.Sprintf("%d", num)

			// ioutil.WriteFile(outputFilePath+nameext+testNum, data, 0644)
		}
	}
	return nil
}

// func main() {
// 	inputFilePath := "resultsMA.json"
// 	outputFilePath := "output.json"

// 	results, err := readPostmanResultsFromFile(inputFilePath)
// 	if err != nil {
// 		fmt.Println("Error reading input file:", err)
// 		return
// 	}

// 	// fmt.Println("Data extracted from", results)

// 	err = writeExtractedDataToFile(results, outputFilePath)
// 	if err != nil {
// 		fmt.Println("Error writing output file:", err)
// 		return
// 	}

// 	fmt.Println("Data extracted and written to", outputFilePath)
// }
