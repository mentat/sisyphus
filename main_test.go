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
	JobStatus = NewThreadSafeMap()

	rid, err := generateUniqueID()

	if err != nil {
		t.Fatalf(err.Error())
	}

	payload := url.Values{
		"apiUser":      {TestUserName},
		"apiKey":       {TestApiKey},
		"sourceDC":     {"IAD"},
		"sourceBucket": {"testing"},
		"sourceFile":   {"pic.jpg"},
		"destDC":       {"DFW"},
		"destBucket":   {"testing"},
		"destFile":     {"pic.jpg"},
		"requestID":    {rid},
	}

	req, _ := http.NewRequest("POST", "/api/copy", strings.NewReader(payload.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Error copying file: %v", w.Code)
	}
}

func TestApiHandlerAsync(t *testing.T) {
	/*
	   Ensure API Handler can read HTTP POST to route.
	*/
	JobStatus = NewThreadSafeMap()

	fmt.Println("Test API async handler...")
	rid, err := generateUniqueID()

	if err != nil {
		t.Fatalf(err.Error())
	}

	payload := url.Values{
		"apiUser":      {TestUserName},
		"apiKey":       {TestApiKey},
		"sourceDC":     {"IAD"},
		"sourceBucket": {"testing"},
		"sourceFile":   {"pic.jpg"},
		"destDC":       {"DFW"},
		"destBucket":   {"testing"},
		"destFile":     {"pic.jpg"},
		"async":        {"1"},
		"requestID":    {rid},
	}

	req, _ := http.NewRequest("POST", "/api/copy", strings.NewReader(payload.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	Router().ServeHTTP(w, req)

	if w.Code != 202 {
		t.Fatalf("Error copying file: %v", w.Code)
	}

	for {
		getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/copy/%s", rid), nil)

		wtr := httptest.NewRecorder()

		Router().ServeHTTP(wtr, getReq)

		if wtr.Code == 201 {
			// Request complete
			break
		} else if wtr.Code == 503 {
			t.Fatalf("Request copy failed.")
		} else {
			// Waiting for request...
		}
	}
}

func TestApiHandlerToken(t *testing.T) {
	/*
	   Ensure API Handler can read HTTP POST to route.
	*/
}
