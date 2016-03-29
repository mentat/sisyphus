package main

import "testing"

func TestUID(t *testing.T) {

	uid, err := generateUniqueID()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if uid == "" {
		t.Fatalf("UID is empty.")
	}
}
