package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
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

func handle(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		fmt.Println("connect: ", r.Proto)

		if pusher, ok := w.(http.Pusher); ok {
			err := pusher.Push("/main.go", nil)
			if err != nil {
				fmt.Println("push err:", err)
			}
		}

		//w.Write([]byte("Hello"))
	}
	if r.RequestURI == "/test_push" {
		fmt.Println("get push")
		w.Write([]byte("get push"))
	}

}