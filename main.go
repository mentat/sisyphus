package main

import (
	"fmt"
	"gocloudfiles"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var JobStatus *ThreadSafeMap

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/copy", CopyCloudFile)
	r.HandleFunc("/api/copy/{requestID}", GetCopyStatus)
	return r
}

func main() {

	JobStatus = NewThreadSafeMap()

	InitLogs(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	http.Handle("/", Router())
	http.ListenAndServe(":12345", nil)
}

func runCopy(cf *gocloudfiles.CloudFiles, requestID, sourceDC, sourceBucket,
	sourceFile, destDC, destBucket, destFile string) error {

	stat := Status{
		Done:    false,
		Running: true,
	}

	JobStatus.Set(requestID, stat)

	err := cf.CopyFile(sourceDC, sourceBucket, sourceFile, destDC, destBucket,
		destFile)

	stat.Done = true
	stat.Running = false

	if err != nil {
		stat.Error = err.Error()
	}

	JobStatus.Set(requestID, stat)

	return err
}

// CopyCloudFile ...
func CopyCloudFile(w http.ResponseWriter, r *http.Request) {
	/*
		Copy a cloud file from one DC to another.
	*/
	requestID := r.FormValue("requestID")
	apiUser := r.FormValue("apiUser")
	apiKey := r.FormValue("apiKey")
	token := r.FormValue("token")
	sourceDC := r.FormValue("sourceDC")
	sourceBucket := r.FormValue("sourceBucket")
	sourceFile := r.FormValue("sourceFile")
	destDC := r.FormValue("destDC")
	destBucket := r.FormValue("destBucket")
	destFile := r.FormValue("destFile")
	asyncMode := r.FormValue("async")

	if requestID == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, Response{
			"status":  400,
			"message": "Request ID required.",
		})
		return
	}

	var cf *gocloudfiles.CloudFiles

	if apiUser == "" && token != "" {
		cf = gocloudfiles.NewCloudFilesImpersonation(token)
		if err := cf.RefreshCatalog(); err != nil {
			w.WriteHeader(400)
			fmt.Println(err.Error())
			fmt.Fprint(w, Response{"status": 400, "message": err.Error()})
			return
		}
	} else {
		cf = gocloudfiles.NewCloudFiles(apiUser, apiKey)
		if err := cf.Authorize(); err != nil {
			w.WriteHeader(400)
			fmt.Println(err.Error())
			fmt.Fprint(w, Response{"status": 400, "message": err.Error()})
			return
		}
	}

	if asyncMode == "1" {
		go runCopy(cf, requestID, sourceDC, sourceBucket, sourceFile, destDC,
			destBucket, destFile)

		w.WriteHeader(202)
		fmt.Fprint(w, Response{
			"status":     202,
			"message":    "The copy request is pending.",
			"request_id": requestID,
		})

	} else {
		err := runCopy(cf, requestID, sourceDC, sourceBucket, sourceFile,
			destDC, destBucket, destFile)

		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err.Error())
			fmt.Fprint(w, Response{"status": 500, "message": err.Error()})
			return
		}

		fmt.Fprint(w, Response{
			"status":     201,
			"message":    "Request has completed.",
			"request_id": requestID,
		})
	}
}

// GetCopyStatus ...
func GetCopyStatus(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	requestID := vars["requestID"]

	stat, ok := JobStatus.Get(requestID)
	if !ok {
		w.WriteHeader(404)
		fmt.Fprint(w, Response{
			"status":     404,
			"message":    "Cannot find request.",
			"request_id": requestID,
		})
		return
	}

	if stat.Error != "" {
		w.WriteHeader(503)
		fmt.Fprint(w, Response{
			"status":     503,
			"message":    stat.Error,
			"request_id": requestID,
		})
	} else {
		if stat.Running {
			w.WriteHeader(200) // Accepted
			fmt.Fprint(w, Response{
				"status":     200,
				"message":    "Request is still running.",
				"request_id": requestID,
			})
		} else if stat.Done {
			w.WriteHeader(201) // Request fullfilled
			fmt.Fprint(w, Response{
				"status":     201,
				"message":    "Request has completed.",
				"request_id": requestID,
			})
		} else {
			w.WriteHeader(202) // Request pending
			fmt.Fprint(w, Response{
				"status":     202,
				"message":    "The copy request is pending.",
				"request_id": requestID,
			})
		}
	}
}
