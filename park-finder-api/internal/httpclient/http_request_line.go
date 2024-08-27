package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (h HTTPClient) SendLinePayRequestAPI(req models.LineReserveAPIRequest, url string, res *models.LineReserveAPIResponse) error {
	header := utility.GenHeaderLinePay(url, req)

	requestBody, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return err
	}

	request, err := http.NewRequest(http.MethodPost, os.Getenv("REQUEST_LINE_RESERVE"), strings.NewReader(string(requestBody)))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return err
	}

	for key, value := range header {
		request.Header.Add(key, value)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
		return err
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Printf("%s\n", requestBody)
	fmt.Println("Send Line Pay API Reserve Status code:", response.StatusCode)

	err = json.Unmarshal(responseBody, res)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	if res.ReturnCode != "0000" {
		return errors.New(res.ReturnMessage)
	}

	return nil

}

func (h HTTPClient) SendLinePayConfirmAPI(req models.LineConfirmAPIRequest, url, tran string, res *models.LineConfirmAPIResponse) error {
	header := utility.GenHeaderLinePay(url, req)

	requestBody, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return err
	}

	request, err := http.NewRequest(http.MethodPost, "https://api-pay.line.me/v3/payments/"+tran+"/confirm", strings.NewReader(string(requestBody)))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return err
	}

	for key, value := range header {
		request.Header.Add(key, value)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Error pretty-printing JSON:", err)
		return err
	}

	fmt.Println("Response Body:")
	fmt.Println(prettyJSON.String())
	fmt.Printf("%s\n", requestBody)
	fmt.Println("Send Line Pay API Reserve Status code:", response.StatusCode)

	err = json.Unmarshal(responseBody, res)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}
	return nil

}
