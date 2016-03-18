package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestApiHandler(t *testing.T) {
	/*
	   Ensure API Handler can read HTTP POST to route.
	*/
	fmt.Println("Test API handler...")
	payload := url.Values{
		"apiUser":      {TestUserName},
		"apiKey":       {TestApiKey},
		"sourceDC":     {"IAD"},
		"sourceBucket": {"testing"},
		"sourceFile":   {"pic.jpg"},
		"destDC":       {"DFW"},
		"destBucket":   {"testing"},
		"destFile":     {"pic.jpg"},
	}

	req, _ := http.NewRequest("POST", "/api/copy", strings.NewReader(payload.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	http.HandlerFunc(CopyCloudFile).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Error copying file: %v", w.Code)
	}
}
