package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello~~~~~~~~")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":9487", nil)
}

func main() {
	handleRequests()
}
