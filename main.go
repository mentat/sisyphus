package main

import (
	"fmt"
	"gocloudfiles"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	InitLogs(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	r := mux.NewRouter()
	r.HandleFunc("/api/copy", CopyCloudFile)

	http.ListenAndServe(":12345", r)

}

func CopyCloudFile(w http.ResponseWriter, r *http.Request) {
	/*
		Copy a cloud file from one DC to another.
	*/
	apiUser := r.FormValue("apiUser")
	apiKey := r.FormValue("apiKey")
	sourceDC := r.FormValue("sourceDC")
	sourceBucket := r.FormValue("sourceBucket")
	sourceFile := r.FormValue("sourceFile")
	destDC := r.FormValue("destDC")
	destBucket := r.FormValue("destBucket")
	destFile := r.FormValue("destFile")

	cf := gocloudfiles.NewCloudFiles(apiUser, apiKey)
	if err := cf.Authorize(); err != nil {
		w.WriteHeader(400)
		fmt.Println(err.Error())
		fmt.Fprint(w, Response{"status": 400, "message": err.Error()})
		return
	}

	err := cf.CopyFile(sourceDC, sourceBucket, sourceFile, destDC, destBucket, destFile)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		fmt.Fprint(w, Response{"status": 500, "message": err.Error()})
		return
	}

	fmt.Fprint(w, Response{"status": 200,
		"message": "File successfully copied."})
}
