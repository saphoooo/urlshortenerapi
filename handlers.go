package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/catinello/base62"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	longURL := string(body)
	shortURL := base62.Encode(strToInt(longURL))
	log.Printf("New short URL %s, for %s...", shortURL, longURL)

	conn := pool.Get()
	defer conn.Close()

	// Create a new entry in Redis only if the key (base62 URL encoded) doesn't exist
	exists, err := redis.Int(conn.Do("EXISTS", shortURL))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if exists == 0 {
		_, err = conn.Do("SET", shortURL, longURL)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

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
