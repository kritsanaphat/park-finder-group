package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (h HTTPClient) SendCancelReserveAPI(order_id string) {
	url := fmt.Sprintf("%s/cancel_reserve?order_id=%s", os.Getenv("CRONJOB_HOST"), order_id)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Println("Send Cancel Reserve API Reserve Status code:", response.StatusCode)
}

func (h HTTPClient) SendCancelExtendReserveAPI(order_id string, hour_end int) {
	url := fmt.Sprintf("%s/cancel_extend_reserve?order_id=%s", os.Getenv("CRONJOB_HOST"), order_id)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Println("Send Cancel Reserve API Reserve Status code:", response.StatusCode)
}

func (h HTTPClient) SendCancelReserveInAdvanceAPI(order_id, email string) {
	url := fmt.Sprintf("%s/cancel_reserve_in_advance?order_id=%s", os.Getenv("CRONJOB_HOST"), order_id+","+email)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Println("Send Cancel Reserve In Advance API Reserve Status code:", response.StatusCode)
}

func (h HTTPClient) SendTimeoutReserveAPI(order_id, date, hour, min string) {
	url := fmt.Sprintf("%s/timeout_reserve?order_id=%s&date=%s&hour=%s&min=%s", os.Getenv("CRONJOB_HOST"), order_id, date, hour, min)
	fmt.Println(url)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Println("Send Timeout Reserve API Reserve Status code:", response.StatusCode)
}

func (h HTTPClient) SendOpenStatusArea(parking_id string, range_time int) {
	url := fmt.Sprintf("%s/update_open_area_status?parking_id=%s&range_time=%d", os.Getenv("CRONJOB_HOST"), parking_id, range_time)
	fmt.Println(url)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Println("Send Send open status area API Reserve Status code:", response.StatusCode)
}

func (h HTTPClient) SendRemoveJobAPI(jobId string) {
	url := fmt.Sprintf("%s/remove_job?job_id=%s", os.Getenv("CRONJOB_HOST"), jobId)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Println("Send Remove Job Status code:", response.StatusCode)
}
