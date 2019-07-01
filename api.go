package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/no-of-pods", getNoOfPods)
	http.HandleFunc("/start-app", startAppHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
