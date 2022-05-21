package server

import (
	"encoding/json"
	"fmt"
	"gorpc/service"
	"gorpc/serviceHandler"
	"log"
	"net/http"
	"path/filepath"
)

type HttpHandler struct {
	handler *serviceHandler.ServiceHandler
}

type HttpServiceBody struct {
	Encoding int `json:"encoding"`
	Compress int `json:"compress"`
	Args []byte `json:"args"`
	Timeout int `json:"timeout"`
	Service string `json:"service"`
}

func ListenHTTP(s *serviceHandler.ServiceHandler) {
	srv := &http.Server{Addr: ":8000"}
	http.Handle("/", HttpHandler{handler: s})
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
}

func (this HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	serviceBody := new(HttpServiceBody)
	decoder.Decode(serviceBody)

	serviceMethod := serviceBody.Service
	var method *service.MethodType
	fmt.Println("method:", serviceMethod)
	err, method := this.handler.ResolveServiceMethod(serviceMethod)
	if err != nil {
		fmt.Println("ResolveServiceMethod err:", err)
		w.Write([]byte("service not found"))
		return
	}

	args := method.UnmarshalArgs([]byte{8,1,16,2})

	reply := method.NewReply()

	method.Call(args, reply)
	encoder := json.NewEncoder(w)
	encoder.Encode(reply)
}