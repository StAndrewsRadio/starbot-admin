package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Reads from a URL.
func ReadFromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := resp.Body.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n", html), nil
}
