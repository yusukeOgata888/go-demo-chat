package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

// テンプレートを表す。
type templateHandler struct {
	once  sync.Once
	filename string
	templ *template.Template
}

// ServeHTTPはHTTPリクエストを処理します。
func (t *templateHandler)  ServeHTTP(w http.ResponseWriter, r *http.Request){
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}
func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("OgaYu@5247")
	gomniauth.WithProviders(
		facebook.New("265709657757-jogbr1g7fmvhhdsvbgjdfg4jncejo87e.apps.googleusercontent.com","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", "http://localhost:8080/auth/callback/facebook"),
		github.New("265709657757-jogbr1g7fmvhhdsvbgjdfg4jncejo87e.apps.googleusercontent.com","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","http://localhost:8080/auth/callback/github"),
		google.New("978332130435-5r3irau9dpblnhtjl90rev51l0e1qdkf.apps.googleusercontent.com","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","http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/",loginHandler)
	http.Handle("/room", r)
	// チャットルームを開始します。
	go r.run()
	// webサーバーを起動します。
	log.Println("Webサーバーを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}