package paricheh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ParichehBaseURL = "https://paricheh.roboepics.com"
)

type SubmissionRequest struct {
	Flag string `json:"flag"`
}

func SendSubmittedFlag(problemPath, token, flag string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("invalid AccessToken")
	}

	fullPath := fmt.Sprintf("%s/submission/%s", ParichehBaseURL, problemPath)

	buffer := &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(&SubmissionRequest{
		Flag: flag,
	}); err != nil {
		return false, fmt.Errorf("failed to encode request body")
	}

	req, err := http.NewRequest("POST", fullPath, buffer)
	if err != nil {
		return false, fmt.Errorf("failed to instantiate request")
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error while requesting server: %v", err)
	}

	switch response.StatusCode {
	case http.StatusTeapot:
		return false, nil
	case http.StatusOK:
		return true, nil
	}

	responseBytes, _ := ioutil.ReadAll(response.Body)

	return false, fmt.Errorf("invalid response from server: %d %q", response.StatusCode, string(responseBytes))
}
