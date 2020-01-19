package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/catinello/base62"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

const base = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type url struct {
	URL string `json:"url,omitempty"`
}

func strToInt(str string) int {
	res := 0
	for _, r := range str {
		res = (62 * res) + strings.Index(base, string(r))
	}
	if res < 0 {
		return -res
	}
	return res
}

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

	url := &url{URL: shortURL}
	resp, err := json.Marshal(url)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", resp)
}

func originalHandler(w http.ResponseWriter, r *http.Request) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	v := mux.Vars(r)

	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Int(conn.Do("EXISTS", v["url"]))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if exists == 0 {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	u, err := redis.String(conn.Do("GET", v["url"]))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	originalurl := &url{URL: u}
	resp, err := json.Marshal(originalurl)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", resp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/new", shortenerHandler).Methods("POST")
	r.HandleFunc("/api/v1/{url}", originalHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
