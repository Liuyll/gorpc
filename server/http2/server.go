package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	srv := &http.Server{Addr: ":8000"}
	http.Handle("/", handle)
	certPath, err := filepath.Abs("D:\\code\\utils\\https\\server.crt")
	if err != nil {
		log.Fatal(err)
	}
	keyPath, err := filepath.Abs("D:\\code\\utils\\https\\server_no_passwd.key")
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Duration(100) * time.Second)
}

type HttpServiceBody struct {
	Encoding int `json:"encoding"`
	Compress int `json:"compress"`
	Args interface{} `json:"args"`
	Timeout int `json:"timeout"`
}

func ListenHTTP() {
	srv := &http.Server{Addr: ":8000"}
	http.Handle("/", handle)
	certPath, err := filepath.Abs("D:\\code\\utils\\https\\server.crt")
	if err != nil {
		log.Fatal(err)
	}
	keyPath, err := filepath.Abs("D:\\code\\utils\\https\\server_no_passwd.key")
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Duration(100) * time.Second)
}

func handle(w http.ResponseWriter, r *http.Request, ) {
	// delete start /
	url := r.RequestURI[1:]

	m := strings.Split(url, "/")
	if len(m) != 2 {
		return
	}

	serviceName, methodName := m[0], m[1]
	decoder := json.NewDecoder(r.Body)

	serviceBody := new(HttpServiceBody)
	decoder.Decode(serviceBody)
}