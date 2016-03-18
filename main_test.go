package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

var (
	TestUserName = os.Getenv("TEST_USERNAME")
	TestApiKey   = os.Getenv("TEST_KEY")
)

func TestMain(m *testing.M) {
	if TestUserName == "" || TestApiKey == "" {
		fmt.Println("Please set the environment variables TEST_USERNAME and TEST_KEY")
		os.Exit(1)
	} else {
		os.Exit(m.Run())
	}
}

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
