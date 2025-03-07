package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ProvarResults2Splunk(excelFile, splunkHost, splunkPort, userName, password, source, index string) error {

	splunkUrl := "https://" + splunkHost + ":" + splunkPort
	jsonArrayFile, err := Excel2Json(excelFile)
	if err != nil {
		// fmt.Println(err)
		return err
	}

	userData := []byte(userName + ":" + password)
	basicAuth := base64.StdEncoding.EncodeToString(userData)
	// fmt.Println("Basic Auth: ", basicAuth)

	var data []map[string]interface{}
	var jsonData map[string]interface{}

	jsonFile, err := ioutil.ReadFile(jsonArrayFile)
	if err != nil {
		// fmt.Println(err)
		return err
	}

	// Unmarshal the JSON array into a slice of maps
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Loop through the slice and marshal each map to JSON
	for cnt, obj := range data {
		jsonObj, err := json.Marshal(obj)
		if err != nil {
			// fmt.Println(err)
			return err
		}
		// Print each JSON object as a string
		// fmt.Println(string(jsonObj))

		err = json.Unmarshal(jsonObj, &jsonData)
		if err != nil {
			// fmt.Println(err)
			return err
		}

		testPath := fmt.Sprintf("%v", jsonData["testPath"])

		var pathList []string
		if strings.Contains(testPath, "\\") {
			pathList = strings.Split(testPath, "\\")
		} else if strings.Contains(testPath, "/") {
			pathList = strings.Split(testPath, "/")
		}

		// fmt.Println("TestType:", string(os.PathSeparator), jsonData["TestPath"])
		// fmt.Println("TestPath:", testPath, pathList[2])
		jsonData["testType"] = pathList[2]
		jsonData["functionality"] = pathList[3]
		// Add more fields to the map
		jsonData["jobName"] = os.Getenv("BUILD_DEFINITIONNAME")
		jsonData["runID"] = os.Getenv("BUILD_BUILDNUMBER")
		// jsonData["releaseNo"] = os.Getenv("ReleaseNo")
		// jsonData["environment"] = os.Getenv("Environment")
		// jsonData["productTeam"] = os.Getenv("productTeam")

		// Marshal the map back to JSON
		jsonObject, err := json.Marshal(jsonData)
		if err != nil {
			// fmt.Println(err)
			return err
		}
		fmt.Printf("Test Result %v ", cnt+1)
		// fmt.Printf("Test Result %v : %v \n", cnt+1, string(jsonObject))
		PostResults(splunkUrl, source, index, basicAuth, jsonObject)
		// Print the new JSON object

	}
	return nil
}

func Json2Splunk(jsonArrayFile, splunkHost, splunkPort, userName, password, source, index string) error {

	splunkUrl := "https://" + splunkHost + ":" + splunkPort

	userData := []byte(userName + ":" + password)
	basicAuth := base64.StdEncoding.EncodeToString(userData)
	// fmt.Println("Basic Auth: ", basicAuth)

	var data []map[string]interface{}
	var jsonData map[string]interface{}

	jsonFile, err := ioutil.ReadFile(jsonArrayFile)
	if err != nil {
		// fmt.Println(err)
		return err
	}

	// Unmarshal the JSON array into a slice of maps
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Loop through the slice and marshal each map to JSON
	for cnt, obj := range data {
		jsonObj, err := json.Marshal(obj)
		if err != nil {
			// fmt.Println(err)
			return err
		}
		// Print each JSON object as a string
		// fmt.Println(string(jsonObj))

		err = json.Unmarshal(jsonObj, &jsonData)
		if err != nil {
			// fmt.Println(err)
			return err
		}

		currentTime := time.Now()

		jsonData["testDate"] = currentTime.Format("2006.01.02 15:04:05")

		// Marshal the map back to JSON
		jsonObject, err := json.Marshal(jsonData)
		if err != nil {
			// fmt.Println(err)
			return err
		}
		fmt.Printf("Test Result %v ", cnt+1)
		// fmt.Printf("Test Result %v : %v \n", cnt+1, string(jsonObject))
		PostResults(splunkUrl, source, index, basicAuth, jsonObject)
		// Print the new JSON object

	}
	return nil
}

// PostResults : to post the results to Splunk Dashboards
func PostResults(splunkHost, source, index, auth string, postReq []byte) {

	splunkURL := splunkHost + "/services/receivers/stream?sourcetype=" + source + "&index=" + index
	req, err := http.NewRequest("POST", splunkURL, bytes.NewBuffer(postReq))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// fmt.Println("Splunk Request Body:", bytes.NewBuffer(postReq))

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("The HTTP request failed with error ", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	// fmt.Println("splunkURL", splunkURL)

}

// function to read the JMeter csv results file and convert to JSON format (using the CSV2Json funtion) and added few more fileds to the JSON object with Jenkins Build Number and Job and push the results to Splunk
func JMeterResults2Splunk(jmFile, splunkHost, splunkPort, userName, password, source, index string) error {
	splunkUrl := "https://" + splunkHost + ":" + splunkPort
	userData := []byte(userName + ":" + password)
	basicAuth := base64.StdEncoding.EncodeToString(userData)

	jsonArrayFile, err := CSV2Json(jmFile)
	if err != nil {
		fmt.Println("error converting to json", err)
		return err
	}

	var data []map[string]interface{}
	var jsonData map[string]interface{}

	jsonFile, err := ioutil.ReadFile(jsonArrayFile)
	if err != nil {
		fmt.Println("Reading the Json file", err)
		return err
	}

	// Unmarshal the JSON array into a slice of maps
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Loop through the slice and marshal each map to JSON
	for cnt, obj := range data {
		jsonObj, err := json.Marshal(obj)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// Print each JSON object as a string
		// fmt.Println(string(jsonObj))

		err = json.Unmarshal(jsonObj, &jsonData)
		if err != nil {
			fmt.Println(err)
			return err
		}

		req := strings.Split(fmt.Sprintf("%v", jsonData["threadName"]), " ")
		testSuite := ""
		for t := 0; t < len(req)-1; t++ {
			if t < 1 {
				testSuite = testSuite + "" + req[t]
			} else {
				testSuite = testSuite + " " + req[t]
			}
		}
		i, err := strconv.ParseInt(fmt.Sprintf("%v", jsonData["timeStamp"]), 10, 64)
		if err != nil {
			panic(err)
		}
		tm := time.UnixMilli(i)

		jsonData["testSuite"] = testSuite
		jsonData["testName"] = fmt.Sprintf("%v", jsonData["label"])
		jsonData["status"] = fmt.Sprintf("%v", jsonData["success"])
		jsonData["testURL"] = fmt.Sprintf("%v", jsonData["URL"])
		jsonData["testTime"] = tm.Format("15:04:05.000")
		jsonData["testDate"] = tm.Format("2006-01-02")

		jsonData["executionTime"] = fmt.Sprintf("%v", jsonData["Latency"])
		jsonData["reason"] = fmt.Sprintf("%v", jsonData["failureMessage"])
		jsonData["httpResponse"] = fmt.Sprintf("%v", jsonData["responseMessage"])
		jsonData["statusCode"] = fmt.Sprintf("%v", jsonData["responseCode"])
		jsonData["automationTool"] = "JMeter"

		// Add more fields to the map
		if os.Getenv("BUILD_NUMBER") != "" {
			jsonData["jobName"] = os.Getenv("JOB_NAME")
			jsonData["runID"] = os.Getenv("BUILD_NUMBER")
		}

		if os.Getenv("BUILD_BUILDNUMBER") != "" {
			jsonData["jobName"] = os.Getenv("BUILD_DEFINITIONNAME")
			jsonData["runID"] = os.Getenv("BUILD_BUILDNUMBER")
		}

		jsonData["releaseNo"] = os.Getenv("ReleaseNo")
		jsonData["environment"] = os.Getenv("Environment")
		jsonData["productTeam"] = os.Getenv("productTeam")

		jsonData["Tenant"] = os.Getenv("Tenant")
		jsonData["Priority"] = os.Getenv("Priority")
		jsonData["Microsite"] = os.Getenv("Microsite")

		delete(jsonData, "Connect")
		delete(jsonData, "IdleTime")
		delete(jsonData, "Latency")
		delete(jsonData, "dataType")
		delete(jsonData, "elapsed")
		delete(jsonData, "grpThreads")
		delete(jsonData, "allThreads")
		delete(jsonData, "URL")

		// Marshal the map back to JSON
		jsonObject, err := json.Marshal(jsonData)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// fmt.Printf("Test Result %v ", cnt+1)
		fmt.Printf("Test Result %v : %v \n", cnt+1, string(jsonObject))
		PostResults(splunkUrl, source, index, basicAuth, jsonObject)
		// Print the new JSON object

	}
	return nil
}
