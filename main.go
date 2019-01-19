package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Server settings
const (
	isLogInFile  = false
	isDemoUser   = true
	isDemoSignIn = false && isDemoUser
)

func main() {
	if isLogInFile {
		f, error := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if error != nil {
			log.Fatal(error)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ptrn_register", registerHandler)
	http.HandleFunc("/ptrn_registerRouter", registerRouterHandler)
	http.HandleFunc("/ptrn_login", loginHandler)
	http.HandleFunc("/ptrn_loginRouter", loginRouterHandler)
	http.HandleFunc("/ptrn_auth", authHandler)
	http.HandleFunc("/ptrn_signoutRouter", signoutRouterHandler)

	s := &http.Server{
		Addr:           ":8088",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}
	log.Println("main", "SERVER STARTED")
	s.ListenAndServe()
}
