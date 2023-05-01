package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	Port  int    `ini:"port" cfg:"port" cfgDefault:"8080" cfgHelper:"Port"`
	DBURL string `ini:"dburl" cfg:"dburl" cfgHelper:"Database URL" cfgRequired:"true"`
}

var (
	//go:embed assets/*
	assets embed.FS
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()
		fmt.Println(r.Form.Get("message"))
	}

	index, err := assets.ReadFile("assets/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("index.html").Parse(string(index))
	if err != nil {
		log.Fatal(err)
	}

	data := struct{}{}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	cfg := Config{}

	config.File = "forum.ini"
	err := config.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/forum", homeHandler)

	s := &http.Server{
		Handler:        mux,
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on port %d\n", cfg.Port)
	log.Fatal(s.ListenAndServe())

}
