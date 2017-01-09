package main

import (
	"fmt"
	"github.com/dgruber/gonfig"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"
)

var config = struct {
	sync.RWMutex
	data map[string]interface{}
}{data: make(map[string]interface{})}

func updateConfig() {
	for range time.NewTicker(5 * time.Second).C {
		conf, err := gonfig.FetchConfig()
		config.Lock()
		config.data = conf
		config.Unlock()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("updated configuration")
		}
	}
}

func init() {
	go updateConfig()
}

func renderIndex(w http.ResponseWriter) error {
	tplt, err := template.ParseFiles("./templates/index.template")
	if err != nil {
		return err
	}
	config.RLock()
	errExec := tplt.Execute(w, config.data)
	config.RUnlock()
	return errExec
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	if err := renderIndex(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = ":8888"
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	//http.Handle("/", r)

	http.ListenAndServe(port, r)
}
