package main

import (
	"authMicroservice/internal/handlers"
	"authMicroservice/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	Key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

func main() {
	//Configuration authentication
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	goth.UseProviders(google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:8080/auth/google/callback"))

	store := sessions.NewCookieStore([]byte(Key))
	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd
	gothic.Store = store

	m := map[string]string{
		"google": "Google",
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex := &models.ProviderIndex{Providers: keys, ProvidersMap: m}

	//Define server
	p := pat.New()

	//Login
	p.Get("/auth/{provider}/callback", handlers.CallbackHandler)
	p.Get("/auth/{provider}", handlers.LoginHandler)
	p.Get("/logout/{provider}", handlers.LogoutHandler)

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		err := t.Execute(res, providerIndex)
		if err != nil {
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", p))
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`
