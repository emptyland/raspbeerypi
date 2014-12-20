package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"api"
	"model"
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

// api/job/list
// api/job/content
// api/job/run
func serve(conf *Conf) {
	handler := http.NewServeMux()

	handler.Handle("/api/state", api.NewService(&model.StateModel{}))
	handler.Handle("/api/memory", api.NewService(&model.MemoryModel{}))
	handler.Handle("/api/disk", api.NewService(&model.DiskUsageModel{}))
	handler.Handle("/api/job/", model.NewJobService())

	proxy := &fileServerProxy{
		Handler: http.FileServer(http.Dir(conf.HTDoc.Root)),
	}
	handler.Handle("/", proxy)

	server := &http.Server{
		Addr:         conf.Address,
		Handler:      handler,
		ReadTimeout:  time.Duration(conf.Timeout) * time.Second,
		WriteTimeout: time.Duration(conf.Timeout) * time.Second,
	}
	log.Fatal(server.ListenAndServeTLS(conf.PemCrt, conf.PemKey))
}

type fileServerProxy struct {
	Handler http.Handler
}

var _ = (http.Handler)(&fileServerProxy{})

const (
	kNavTab        = "nav-tab"
	kNavTabDefault = "dashboard"
)

func (self *fileServerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// proxy it!
		self.Handler.ServeHTTP(w, r)
		return
	}

	cookie, err := r.Cookie(kNavTab)
	if err != nil {
		cookie = &http.Cookie{
			Name:  kNavTab,
			Value: kNavTabDefault,
			//Domain: "/",
			Expires: time.Now().Add(2 * time.Hour),
		}
	} else {
		log.Println("cookie:", cookie)
	}

	http.SetCookie(w, cookie)
	// proxy it!
	self.Handler.ServeHTTP(w, r)
}
