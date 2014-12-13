package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"api"
)

var _ = fmt.Println

type Conf struct {
	PemCrt  string   `json:"pem_crt"`
	PemKey  string   `json:"pem_key"`
	Address string   `json:"address"`
	Timeout int64    `json:"timeout"`
	HTDoc   HTDocDef `json:"htdoc"`
}

type HTDocDef struct {
	Root    string            `json:"root"`
	Mapping map[string]string `json:"mapping"`
}

func main() {

	fd, err := os.Open("conf.json")
	if err != nil {
		fmt.Println(err)
	}
	defer fd.Close()

	decoder := json.NewDecoder(fd)
	conf := &Conf{}
	decoder.Decode(&conf)

	fmt.Println(conf)
	serve(conf)
}

func serve(conf *Conf) {
	handler := http.NewServeMux()

	handler.Handle("/api/", api.NewService(&api.HelloModel{}))
	handler.Handle("/", http.FileServer(http.Dir(conf.HTDoc.Root)))

	server := &http.Server{
		Addr:         conf.Address,
		Handler:      handler,
		ReadTimeout:  time.Duration(conf.Timeout) * time.Second,
		WriteTimeout: time.Duration(conf.Timeout) * time.Second,
	}
	log.Fatal(server.ListenAndServeTLS(conf.PemCrt, conf.PemKey))
}