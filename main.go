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

var (
	sc *session.Control

	//go:embed assets/*
	assets embed.FS
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("homeHandler")

	sid, sd, ok := sc.Get(r)
	if !ok {
		//http.Redirect(w, r, "/forum/login/", http.StatusFound)
		return
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

	// exec template
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}

	// http.Redirect(w, r, "/payments", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginHandler")
	if r.Method == http.MethodGet {
		index, err := assets.ReadFile("assets/login.html")
		if err != nil {
			log.Fatal(err)
		}
		t, err := template.New("login.html").Parse(string(index))
		if err != nil {
			log.Fatal(err)
		}

		// exec template
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// login logic

	// create session
	sid, sd := sc.Create()

	// save session
	sc.Save(w, sid, sd)

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logoutHandler")
	sid, _, ok := sc.Get(r)
	if !ok {
		//http.Redirect(w, r, "/forum/login", http.StatusFound)
		return
	}

	// remove session
	sc.Delete(w, sid)

	//http.Redirect(w, r, "/forum/login", http.StatusFound)
}

func issueSession() http.Handler {
	fmt.Println("issueSession")
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		//session.Set(sessionUserKey, *githubUser.ID)
		//session.Set(sessionUsername, *githubUser.Login)
		fmt.Println("githubUser id: ", *githubUser.ID)
		fmt.Println("githubUser login: ", *githubUser.Login)
		fmt.Println("githubUser email: ", *githubUser.Email)
		fmt.Println("githubUser name: ", *githubUser.Name)
		fmt.Println("githubUser avatar: ", *githubUser.AvatarURL)
		fmt.Println("githubUser url: ", *githubUser.URL)
		fmt.Println("githubUser html url: ", *githubUser.HTMLURL)
		fmt.Println("githubUser followers: ", *githubUser.Followers)
		fmt.Println("githubUser following: ", *githubUser.Following)
		fmt.Println("githubUser created at: ", *githubUser.CreatedAt)

		//if err := session.Save(w); err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		http.Redirect(w, req, "/forum", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// profileHandler shows a personal profile or a login button (unauthenticated).
func profileHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("profileHandler")
	//	session, err := sessionStore.Get(req, sessionName)
	//	if err != nil {
	//
	// welcome with login button
	//
	//		page, _ := os.ReadFile("home.html")
	//		fmt.Fprint(w, string(page))
	//		return
	//	}
	//
	// authenticated profile
	//
	//	fmt.Fprintf(w, `<p>You are logged in %s!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`, session.Get(sessionUsername))
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
	mux.HandleFunc("/login", loginHandler)
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
