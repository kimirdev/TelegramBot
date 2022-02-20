package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	searchUrlPattern = "https://customsearch.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s"
)

type ItemJson struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type SearchResult struct {
	Items []ItemJson `json:"items"`
}

func GoogleCustomSearchRequest(query string) (SearchResult, error) {
	var sr SearchResult
	query = strings.TrimSpace(query)
	query = url.QueryEscape(query)

	apiKey, exists := os.LookupEnv("GOOGLE_API_KEY")

	if !exists {
		return sr, errors.New("GOOGLE_API_KEY not found in file .env")
	}

	cx, exists := os.LookupEnv("GOOGLE_CX")

	if !exists {
		return sr, errors.New("GOOGLE_CX not found in file .env")
	}

	searchUrl := fmt.Sprintf(searchUrlPattern, apiKey, cx, query)

	resp, err := http.Get(searchUrl)

	if err != nil {
		return sr, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sr, err
	}
	err = json.Unmarshal(body, &sr)

	return sr, err
}
