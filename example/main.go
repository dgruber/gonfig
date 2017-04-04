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

func updateConfig(cfgCh <-chan gonfig.Config) {
	for cfg := range cfgCh {
		config.Lock()
		config.data = cfg
		config.Unlock()
		fmt.Println("updated configuration")
	}
}

func init() {
	cfgCh, _ := gonfig.ConfigChange(time.Second * 1)
	go updateConfig(cfgCh)
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

	http.ListenAndServe(port, r)
}
