package main

import (
	"embed"
	"fmt"
	"html/template"
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

	DB     db.DB
	GitTag string = "dev"
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

	topicID, err := DB.GetTopicID(parameters)
	if err != nil {
		log.Println(err)
	}

	p, err := DB.ListTopics(topicID)
	if err != nil {
		log.Println(err)
	}

	data := struct {
		Topics []db.Topic
	}{
		Topics: p,
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
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

	log.Printf("Forum version %v", GitTag)
	log.Printf("Listening on port %d", config.Port)
	log.Fatal(s.ListenAndServe())

}
