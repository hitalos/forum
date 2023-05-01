package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"crg.eti.br/go/forum/config"
	"crg.eti.br/go/forum/db"
	"crg.eti.br/go/forum/db/pg"
	_ "github.com/lib/pq"
)

var (
	//go:embed assets/*
	assets embed.FS

	DB db.DB
)

func getParameters(prefix string, r *http.Request) ([]string, error) {
	validateSlug := func(s string) bool {
		if s == "" {
			return false
		}

		for _, c := range s {
			if !((c >= 'a' && c <= 'z') ||
				(c >= '0' && c <= '9') ||
				c == '-') {
				return false
			}
		}

		return true
	}

	path := strings.ReplaceAll(r.URL.Path, "//", "/")
	path = strings.TrimPrefix(path, prefix)
	path = strings.TrimSuffix(path, "/")

	if path == "" {
		return []string{}, nil
	}

	a := strings.Split(path, "/")

	for i := range a {
		if !validateSlug(a[i]) {
			return nil, fmt.Errorf("Invalid slug: %q, %q, path: %q", r.URL.Path, a[i], path)
		}
	}

	return a, nil
}

func redirectHandler(prefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, prefix, http.StatusMovedPermanently)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	parameters, err := getParameters("/forum/", r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
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

	for k, v := range parameters {
		fmt.Printf("%d: %s\n", k, v)
	}

	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
	}

	if r.Method == http.MethodPost &&
		r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			return
		}

		for k, v := range r.PostForm {
			fmt.Printf("%s: %s\n", k, v)
		}
	}

	if r.Method == http.MethodPost &&
		r.Header.Get("Content-Type") == "application/json" {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", string(b))
	}

	fmt.Printf("RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Printf("Host: %s\n", r.Host)
	fmt.Printf("RequestURI: %s\n", r.RequestURI)
	fmt.Printf("URL: %s\n", r.URL)
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("Header: %v\n", r.Header)

}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	DB = pg.New()
	err = DB.Open()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/forum", redirectHandler("/forum/"))
	mux.HandleFunc("/forum/", homeHandler)

	s := &http.Server{
		Handler:        mux,
		Addr:           fmt.Sprintf(":%d", config.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on port %d\n", config.Port)
	log.Fatal(s.ListenAndServe())

}
