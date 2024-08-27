package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (h HTTPClient) SendCameraServiceToDectectIncommingCar(module_code, licence_plate string) bool {
	requestBody, err := json.Marshal(map[string]string{
		"module_code":            module_code,
		"customer_license_plate": licence_plate,
	})
	if err != nil {
		fmt.Println("Error marshalling JSON payload:", err)
		return false
	}

	url := fmt.Sprintf("%s/webhook/check_licence_plate", os.Getenv("CAMERA_HOST"))
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return false
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return false
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &jsonResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return false
	}

	if jsonResponse["response"] == "True" {
		return true
	} else {
		return false
	}
}

func (h HTTPClient) SendCameraServiceTCaptureCar(module_code string) (string, error) {
	url := fmt.Sprintf("%s/webhook/capture_camera?module_code=%s", os.Getenv("CAMERA_HOST"), module_code)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return "", err

	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err

	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &jsonResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", err
	}

	if jsonResponse["response"] == nil {
		return "", errors.New("Can't Capture Car")
	}

	responseValue, ok := jsonResponse["response"].(string)
	if !ok {
		return "", errors.New("Response value is not a string")
	}

	return responseValue, nil
}
