package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (h HTTPClient) SendToAutomaticExtendReserve(order_id string) bool {
	url := fmt.Sprintf("http://%s/customer/extend_reserve?order_id=%s&action=automatic", os.Getenv("HOST"), order_id)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return false

	}

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
