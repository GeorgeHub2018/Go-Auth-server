package main

import (
	"errors"
	"log"

	"github.com/rs/xid"
)

type (
	//Bio user bio
	Bio struct {
		Name string
		Age  string
	}
	//User user struct
	User struct {
		Login    string
		Password string
		BIO      Bio
	}
	//Users for users cache
	Users map[string]User

	//Auth auth struct
	Auth struct {
		ID   string
		User User
	}

	//Auths for auths cache
	Auths map[string]Auth
)

var cacheUsers Users
var cacheAuths Auths
var demoAuth Auth

func (a Auths) findByLogin(login string) (au Auth) {
	for _, auth := range cacheAuths {
		if auth.User.Login == login {
			au = auth
			return
		}
	}
	return
}

func (auth Auth) asData() UserData {
	return UserData{
		Name:   auth.User.BIO.Name,
		Age:    auth.User.BIO.Age,
		AuthID: auth.ID}
}

func (c Users) asData() (l ListData) {
	l = ListData{Users: make(map[string]UserData)}
	for i, user := range cacheUsers {
		var authID string
		rangeAuth := cacheAuths.findByLogin(user.Login)
		if &rangeAuth != nil {
			authID = rangeAuth.ID
		}
		l.Users[i] = UserData{Name: user.BIO.Name, Age: user.BIO.Age, AuthID: authID}
	}

	return
}

// generating random guid string
func genXid() string {
	id := xid.New()
	return id.String()
}

func init() {
	cacheUsers = Users{}
	cacheAuths = Auths{}
	if isDemoUser {
		demoUser := User{
			Login:    "admin",
			Password: "admin",
			BIO: Bio{
				Name: "admin",
				Age:  "100500",
			},
		}
		cacheUsers[demoUser.Login] = demoUser
		if isDemoSignIn {
			demoAuth, _ = loginExecute(demoUser.Login)
		}
	}

	log.Println("init", "cacheUsers:", cacheUsers)
	log.Println("init", "cacheAuths:", cacheAuths)
}

func loginCheck(login, pass string) (err error) {
	if user, ok := cacheUsers[login]; ok && user.Password == pass {
		return
	}

	err = errors.New("invalid pass or login")
	log.Println("loginCheck", err.Error())
	return
}

func authCheck(id string) (auth Auth, err error) {
	if auth, ok := cacheAuths[id]; ok {
		return auth, err
	}

	err = errors.New("invalid id")
	log.Println("authCheck", err.Error())
	return
}

func register(login, pass, name, age string) (err error) {
	if _, ok := cacheUsers[login]; ok {
		err = errors.New("user exist")
		log.Println("register", err.Error())
		return
	}

	cacheUsers[login] = User{
		Login:    login,
		Password: pass,
		BIO: Bio{
			Name: name,
			Age:  age,
		},
	}
	log.Println("register", cacheUsers[login])
	return
}

func loginExecute(login string) (auth Auth, err error) {
	if _, ok := cacheUsers[login]; !ok {
		err = errors.New("user not exist")
		log.Println("loginExecute", err.Error())
		return
	}

	out := genXid()
	cacheAuths[out] = Auth{
		ID:   out,
		User: cacheUsers[login],
	}
	log.Println("loginExecute ", cacheAuths)
	auth = cacheAuths[out]
	return
}

func logoutExecute(id string) {
	if _, ok := cacheAuths[id]; !ok {
		err := errors.New("id not found")
		log.Println("logoutExecute", err.Error())
		return
	}

	auth := cacheAuths[id]
	delete(cacheAuths, id)
	log.Println("logoutExecute", auth, cacheAuths)
}
