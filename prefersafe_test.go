package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeRequest(url string, preferSafeValue string) (respBody []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if preferSafeValue != "" {
		req.Header.Add("Prefer", preferSafeValue)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func TestImage(t *testing.T) {
	srv := httptest.NewServer(serveMux())
	defer srv.Close()

	url := srv.URL + "/image.gif"

	// unsafe
	resp, err := makeRequest(url, "")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if len(resp) > 0 {
		t.Errorf("expected zero-length response, got %q", resp)
	}

	// safe
	resp, err = makeRequest(url, "safe")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if !bytes.Equal(resp, emptyGIF) {
		t.Errorf("expected empty GIF, got: %v (%q)", resp, resp)
	}
}

func TestJSONP(t *testing.T) {
	srv := httptest.NewServer(serveMux())
	defer srv.Close()

	url := srv.URL + "/jsonp.js"
	trueResponse := "PreferSafe(true)"
	falseResponse := "PreferSafe(false)"

	// unsafe
	resp, err := makeRequest(url, "")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if !strings.Contains(string(resp), falseResponse) {
		t.Errorf("unexpected response %q", resp)
	}

	// safe
	resp, err = makeRequest(url, "safe")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if !strings.Contains(string(resp), trueResponse) {
		t.Errorf("unexpected response %q", resp)
	}
}
