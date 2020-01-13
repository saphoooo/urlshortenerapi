package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	u, err := url.Parse(string(body))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	longURL := string(body)
	shortURL := base62.Encode(int(binary.BigEndian.Uint64([]byte(u.Host + u.Path))))
	log.Printf("New short URL created %s for %s...", shortURL, longURL)

	conn := pool.Get()
	defer conn.Close()

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

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v\n", shortURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
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
	url, err := redis.String(conn.Do("GET", v["url"]))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("redirecting to %s...\n", url)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/new", shortenerHandler).Methods("POST")
	r.HandleFunc("/api/v1/{url}", redirectHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
