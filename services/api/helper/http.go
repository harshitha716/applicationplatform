package helper

import (
	"io"
	"net/http"

	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func ForwardResponseHeaders(src *http.Response, c *gin.Context) {
	excludedHeaders := map[string]struct{}{
		"Content-Length": {},
		"Date":           {},
	}
	for key, values := range src.Header {
		if _, excluded := excludedHeaders[key]; !excluded {
			for _, value := range values {
				c.Header(key, value)
			}
		}
	}
}

func HttpGet(httpClient HTTPClient, url string, headers map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating GET request: %w", err)
	}

	addHeaders(req, headers)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, resp, nil
}

func HttpPost(httpClient HTTPClient, url string, payload interface{}, headers map[string]string) ([]byte, *http.Response, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating POST request: %w", err)
	}

	addHeaders(req, headers)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, resp, nil
}

func HttpPut(httpClient HTTPClient, url string, body []byte, headers map[string]string) (respBody []byte, err error) {

	buffer := bytes.NewBuffer(body)

	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return nil, fmt.Errorf("Error while creating put request %s", err)

	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err

	}
	resp.Body.Close()

	return respBody, nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
