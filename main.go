package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"

	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
	"crg.eti.br/go/session"
)

type Config struct {
	GithubClientID     string `ini:"github_client_id" cfg:"github_client_id" cfgRequired:"true" cfgHelper:"Github Client ID"`
	GithubClientSecret string `ini:"github_client_secret" cfg:"github_client_secret" cfgRequired:"true" cfgHelper:"Github Client Secret"`
	GithubCallbackURL  string `ini:"github_callback_url" cfg:"github_callback_url" cfgRequired:"true" cfgHelper:"Github Callback URL"`
	DatabaseURL        string `ini:"database_url" cfg:"database_url" cfgRequired:"true" cfgHelper:"Database URL"`
	Port               int    `ini:"port" cfg:"port" cfgDefault:"8080" cfgHelper:"Port"`
}

type sessionData struct {
	OAuthProvider string
	UserName      string
	AvatarURL     string
}

var (
	sc *session.Control

	//go:embed assets/*
	assets embed.FS
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("homeHandler")

	sid, sd, ok := sc.Get(r)
	if !ok {
		fmt.Println("session not found")
		sid, sd = sc.Create()
	}

	// renew session
	sc.Save(w, sid, sd)

	//////////////////////////

	index, err := assets.ReadFile("assets/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("index.html").Parse(string(index))
	if err != nil {
		log.Fatal(err)
	}

	var (
		sdAUX *sessionData
	)

	if sd.Data != nil {
		sdAUX, ok = sd.Data.(*sessionData)
		if !ok {
			log.Fatal("type assertion failed sessionData")
		}
	}
	data := struct {
		SessionData    *sessionData
		GitHubLoginURL string
		LogoutURL      string
	}{
		SessionData:    sdAUX,
		GitHubLoginURL: "/forum/github/login",
		LogoutURL:      "/forum/logout",
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

	// http.Redirect(w, r, "/payments", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logoutHandler")
	sid, _, ok := sc.Get(r)
	if !ok {
		http.Redirect(w, r, "/forum", http.StatusFound)
		return
	}

	// remove session
	sc.Delete(w, sid)

	http.Redirect(w, r, "/forum", http.StatusFound)
}

func issueSession() http.Handler {
	fmt.Println("issueSession")
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sid, sd, ok := sc.Get(r)
		if !ok {
			fmt.Println("session not found")
			sid, sd = sc.Create()
		}

		fmt.Println("githubUser id.........:", *githubUser.ID)
		fmt.Println("githubUser login......:", *githubUser.Login)
		fmt.Println("githubUser email......:", *githubUser.Email)
		fmt.Println("githubUser name.......:", *githubUser.Name)
		fmt.Println("githubUser avatar.....:", *githubUser.AvatarURL)
		fmt.Println("githubUser url........:", *githubUser.URL)
		fmt.Println("githubUser html url...:", *githubUser.HTMLURL)
		fmt.Println("githubUser followers..:", *githubUser.Followers)
		fmt.Println("githubUser following..:", *githubUser.Following)
		fmt.Println("githubUser created at.:", *githubUser.CreatedAt)

		sdAUX := sessionData{
			OAuthProvider: "github",
			UserName:      *githubUser.Name,
			AvatarURL:     *githubUser.AvatarURL,
		}
		sd.Data = &sdAUX

		sc.Save(w, sid, sd)

		http.Redirect(w, r, "/forum", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func main() {
	cfg := Config{}

	config.File = "config.ini"
	err := config.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	const cookieName = "forum_session"
	sc = session.New(cookieName)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			sc.RemoveExpired()
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	//mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.GithubClientID,
		ClientSecret: cfg.GithubClientSecret,
		RedirectURL:  cfg.GithubCallbackURL,
		Endpoint:     githubOAuth2.Endpoint,
	}

	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig

	mux.Handle(
		"/github/login",
		github.StateHandler(
			stateConfig,
			github.LoginHandler(oauth2Config, nil)))
	mux.Handle(
		"/github/callback",
		github.StateHandler(
			stateConfig,
			github.CallbackHandler(oauth2Config, issueSession(), nil)))

	s := &http.Server{
		Handler:        mux,
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	println(cfg.DatabaseURL)

	log.Printf("Listening on port %d\n", cfg.Port)
	log.Fatal(s.ListenAndServe())

}
