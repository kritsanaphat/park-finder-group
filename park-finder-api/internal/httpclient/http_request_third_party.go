package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (h HTTPClient) SendToPTTReward(url string) string {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return ""

	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""

	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &jsonResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ""
	}

	return jsonResponse["barcode_url"].(string)

}
