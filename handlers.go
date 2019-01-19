package main

import (
	"html/template"
	"log"
	"net/http"
)

func templateParseFiles(sitePage string) (*template.Template, error) {
	return template.ParseFiles(
		"templates/index.html",
		"templates/"+sitePage+"/header.html",
		"templates/"+sitePage+"/main.html",
		"templates/footer.html")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("indexHandler", "enter")
	t, err := templateParseFiles("page_index")
	if err != nil {
		log.Println("indexHandler", err.Error())
		return
	}

	t.ExecuteTemplate(w, "index", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("registerHandler", "enter")
	r.ParseForm()
	var errorText string
	if _, ok := r.Form["ErrorText"]; ok {
		errorText = r.Form["ErrorText"][0]
	}

	t, err := templateParseFiles("page_register")
	if err != nil {
		log.Println("registerHandler", err.Error())
		return
	}

	if errorText != "" {
		data := ErrorData{
			ErrorText: errorText,
		}

		t.ExecuteTemplate(w, "index", data)
		return
	}
	t.ExecuteTemplate(w, "index", nil)
}

func registerRouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("registerRouterHandler", "enter")
	r.ParseForm()
	var login, pass, name, age string
	if _, ok := r.Form["login"]; ok {
		login = r.Form["login"][0]
	}
	if _, ok := r.Form["pass"]; ok {
		pass = r.Form["pass"][0]
	}
	if _, ok := r.Form["name"]; ok {
		name = r.Form["name"][0]
	}
	if _, ok := r.Form["age"]; ok {
		age = r.Form["age"][0]
	}

	err := register(login, pass, name, age)
	if err != nil {
		log.Println("registerRouterHandler", err.Error())
		http.Redirect(w, r, "/ptrn_register?ErrorText="+err.Error(), 303)
		return
	}

	auth, err := loginExecute(login)
	if err != nil {
		log.Println("registerRouterHandler", err.Error())
		http.Redirect(w, r, "/ptrn_register?ErrorText="+err.Error(), 303)
		return
	}

	http.Redirect(w, r, "/ptrn_auth?id="+auth.ID, 301)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("loginHandler", "enter")
	r.ParseForm()
	var errorText string
	if _, ok := r.Form["ErrorText"]; ok {
		errorText = r.Form["ErrorText"][0]
	}

	t, err := templateParseFiles("page_login")
	if err != nil {
		log.Println("loginHandler", err.Error())
		return
	}

	if errorText != "" {
		data := ErrorData{
			ErrorText: errorText,
		}

		t.ExecuteTemplate(w, "index", data)
		return
	}

	t.ExecuteTemplate(w, "index", nil)
}

func loginRouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("loginRouterHandler", "enter")
	r.ParseForm()
	var login, pass string
	if _, ok := r.Form["login"]; ok {
		login = r.Form["login"][0]
	}
	if _, ok := r.Form["pass"]; ok {
		pass = r.Form["pass"][0]

	}

	err := loginCheck(login, pass)
	if err != nil {
		log.Println("loginRouterHandler", login, pass, err.Error())
		http.Redirect(w, r, "/ptrn_login?ErrorText="+err.Error(), 303)
		return
	}

	auth, err := loginExecute(login)
	if err != nil {
		log.Println("err ", login, pass, err.Error())
		http.Redirect(w, r, "/ptrn_login?ErrorText="+err.Error(), 303)
		return
	}
	http.Redirect(w, r, "/ptrn_auth?id="+auth.ID, 301)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("authHandler", "enter")
	r.ParseForm()
	var id string
	if _, ok := r.Form["id"]; ok {
		id = r.Form["id"][0]
	}

	log.Println("authHandler", "id", id)
	log.Println("authHandler", "cacheAuths", cacheAuths)
	auth, err := authCheck(id)
	if err != nil {
		log.Println("authHandler", err.Error())
		http.Redirect(w, r, "/ptrn_login?ErrorText="+err.Error(), 303)
		return
	}

	t, err := templateParseFiles("page_auth")
	if err != nil {
		log.Println("authHandler", err.Error())
		http.Redirect(w, r, "/ptrn_login?ErrorText="+err.Error(), 303)
		return
	}

	data := FullUserData{
		User: auth.asData(),
		List: cacheUsers.asData(),
	}
	log.Println("authHandler", data)
	t.ExecuteTemplate(w, "index", data)
}

func signoutRouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("signoutRouterHandler", "enter")
	r.ParseForm()
	var id string
	if _, ok := r.Form["id"]; ok {
		id = r.Form["id"][0]
	}

	logoutExecute(id)
	http.Redirect(w, r, "/", 303)
}
