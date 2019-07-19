package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/no-of-pods", getNoOfPods)
	http.HandleFunc("/start-app", startAppHandler)
	http.HandleFunc("/stop-app", stopAppHandler)
	http.HandleFunc("/delete-app-storage", deleteAppStorageHandler)
	http.HandleFunc("/get-app-status/", getAppStatusHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
