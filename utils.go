package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const IPO_LIST_URL = "https://iporesult.cdsc.com.np/result/companyShares/fileUploaded"
const IPO_CHECK_URL = "https://iporesult.cdsc.com.np/result/result/check"

type IPOList struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Body    []IPOInfo `json:"body"`
}

type IPOInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Scrip          string `json:"scrip"`
	IsFileUploaded bool   `json:"isFileUploaded"`
}

func getIPOList() ([]IPOInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", IPO_LIST_URL, nil)

	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	strJSON, _ := io.ReadAll(resp.Body)

	result := &IPOList{}

	err = json.Unmarshal(strJSON, result)
	if err != nil {
		panic(err)
	}


	return result.Body, nil
}

type AllotmentStatus struct {
	Success bool `json:"success"`
}

type AllotmentRequest struct {
	BOID           string `json:"boid"`
	CompanyShareID int    `json:"companyShareId"`
}

func checkIPO(BOID string, companyShareID int) (bool, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", IPO_CHECK_URL, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	allotmentRequest := AllotmentRequest{
		BOID:           BOID,
		CompanyShareID: companyShareID,
	}

	reqJSON, err := json.Marshal(allotmentRequest)

	req.Body = io.NopCloser(bytes.NewBuffer(reqJSON))

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("%s", resp.Status)
	}

	defer resp.Body.Close()

	strJSON, _ := io.ReadAll(resp.Body)

	result := &AllotmentStatus{}

	err = json.Unmarshal(strJSON, result)

	if err != nil {
		panic(err)
	}

	return result.Success, nil

	return false, nil
}

func printIPOTable(ipoList []IPOInfo) {
	fmt.Printf("%5s\t%40s\t%10s\t%10s\n", "ID", "Name", "Scrip", "IsFileUploaded")
	for _, ipo := range ipoList {
		fmt.Printf("%5d\t%40s\t%10s\t%10t\n", ipo.ID, ipo.Name, ipo.Scrip, ipo.IsFileUploaded)
	}
}
