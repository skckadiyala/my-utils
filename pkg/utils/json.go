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
	"strings"
)

func ProvarResults2Splunk(excelFile, splunkHost, splunkPort, userName, password, source, index string) error {

	splunkUrl := "https://" + splunkHost + ":" + splunkPort
	jsonArrayFile, err := Excel2Json(excelFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	userData := []byte(userName + ":" + password)
	basicAuth := base64.StdEncoding.EncodeToString(userData)
	fmt.Println("Basic Auth: ", basicAuth)

	var data []map[string]interface{}
	var jsonData map[string]interface{}

	jsonFile, err := ioutil.ReadFile(jsonArrayFile)
	if err != nil {
		fmt.Println(err)
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

		testPath := fmt.Sprintf("%v", jsonData["testPath"])
		var pathList []string
		if strings.Contains(testPath, "\\") {
			pathList = strings.Split(testPath, "\\")
		} else if strings.Contains(testPath, "/") {
			pathList = strings.Split(testPath, "/")
		}

		// fmt.Println("TestType:", string(os.PathSeparator), jsonData["TestPath"])

		jsonData["testType"] = pathList[2]
		jsonData["functionality"] = pathList[3]
		// Add more fields to the map
		jsonData["jobName"] = os.Getenv("BUILD_DEFINITIONNAME")
		jsonData["runID"] = os.Getenv("BUILD_BUILDNUMBER")
		jsonData["releaseNo"] = os.Getenv("ReleaseNo")
		jsonData["environment"] = os.Getenv("Environment")
		jsonData["productTeam"] = os.Getenv("productTeam")

		// Marshal the map back to JSON
		jsonObject, err := json.Marshal(jsonData)
		if err != nil {
			fmt.Println(err)
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
	fmt.Println("Splunk Request Body:", bytes.NewBuffer(postReq))

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

}
