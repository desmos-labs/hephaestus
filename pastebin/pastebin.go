package pastebin

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiURL = "https://pastebin.com/api/api_post.php"
)

// Upload uploads the given text using the PasteBin APIs authenticated with the given API key.
// It returns either an error or the paste URL if the data were uploaded successfully.
func Upload(apiKey string, text string) (string, error) {
	data := url.Values{}
	data.Set("api_dev_key", apiKey)
	data.Set("api_option", "paste")
	data.Set("api_paste_code", text)

	reqURL, err := url.ParseRequestURI(apiURL)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, reqURL.String(), strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
