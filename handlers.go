package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/catinello/base62"
	"github.com/gorilla/mux"
)

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	longURL := string(body)
	shortURL := base62.Encode(strToInt(longURL))
	status, err := dbSet(dbDriver, shortURL, longURL)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(status), status)
	}
	log.Printf("New shortened url \"%s\", for original url %s...", shortURL, longURL)

	u := &url{URL: shortURL}
	resp, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", resp)
}

func originalHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u, status, err := dbGet(dbDriver, v["url"])
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(status), status)
	}
	originalurl := &url{URL: u}
	resp, err := json.Marshal(originalurl)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", resp)
}
