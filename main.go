package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Lawry struct {
	Li  string `json:"li"`
	Way string `json:"way"`
}

func myLay(w http.ResponseWriter, r *http.Request) {
	lilies := Lawry{"Hi", "NiHow"}
	json.NewEncoder(w).Encode(lilies)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
	}
	fmt.Fprintf(w, "Host=%q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr=%q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		panic(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q]=%q\n", k, v)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", rootHandler)
	myRouter.HandleFunc("/li", myLay)
	http.ListenAndServe(":9487", myRouter)
}

func main() {
	handleRequests()
}
