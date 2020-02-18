package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Lawry struct {
	Id  int    `json:"id"`
	Li  string `json:"li"`
	Way string `json:"way"`
}

var (
	lilies       []Lawry
	mySigningKey = []byte("hahahahaha")
)

func deleteLay(w http.ResponseWriter, r *http.Request) {
	key, err := strconv.Atoi(getRequestPathVar(r, "id"))
	if err != nil {
		errHandler(err)
	}
	for index, lily := range lilies {
		if lily.Id == key {
			lilies = append(lilies[:index], lilies[index+1:]...)
		}
	}
}

func addLay(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errHandler(err)
	}
	var lily Lawry
	json.Unmarshal(reqBody, &lily)
	lilies = append(lilies, lily)
	fmt.Fprintf(w, "%+v\n added done.", string(reqBody))
	log.Println(lilies)
}

func showLay(w http.ResponseWriter, r *http.Request) {
	key, err := strconv.Atoi(getRequestPathVar(r, "id"))
	if err != nil {
		errHandler(err)
	}
	for _, lily := range lilies {
		if lily.Id == key {
			json.NewEncoder(w).Encode(lily)
			log.Println(lilies)
		}
	}
	log.Println(key)
}

func getRequestPathVar(r *http.Request, data string) string {
	vars := mux.Vars(r)
	key := vars[data]
	log.Println(key)
	return key
}

func errHandler(err error) {
	log.Fatal(err)
	panic(err)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
	}
	fmt.Fprintf(w, "Host=%q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr=%q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		errHandler(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q]=%q\n", k, v)
	}
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Handle("/", isAuthorized(rootHandler))
	myRouter.HandleFunc("/ls/{id}", showLay)
	myRouter.HandleFunc("/add", addLay).Methods("POST")
	myRouter.HandleFunc("/delete/{id}", deleteLay).Methods("DELETE")
	http.ListenAndServe(":9487", myRouter)
}

func main() {
	handleRequests()
}
