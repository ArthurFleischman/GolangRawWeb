package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	Name     string
	password string
}

var users = map[string]user{}

var sessions = map[string]string{}

var list []user

var tmp *template.Template

func init() {
	tmp = template.Must(template.ParseGlob("templates/*.html"))

}

func index(w http.ResponseWriter, r *http.Request) {
	tmp.ExecuteTemplate(w, "index.html", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("username")
		pass := r.FormValue("password")
		u := user{Name: name, password: pass}
		users[name] = u
		logMessage := fmt.Sprintf("user %s registered\n", u.Name)
		log.Printf(logMessage)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmp.ExecuteTemplate(w, "signup.html", nil)
}
func isLogged(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		log.Println("user not logged in")
		return false
	} else if sessions[c.Value] != "" {
		fmt.Println("ok")
		return true
	} else {
		logMessage := fmt.Sprintf("invalid cookie %s\n", c.Value)
		log.Printf(logMessage)
		return false
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//Check if user already has a cookie(aka is logged)
	if isLogged(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		name := r.FormValue("username")
		pass := r.FormValue("password")
		if users[name].Name == name && users[pass].password == pass {
			id := uuid.NewV4()
			c := &http.Cookie{
				Name:  "session",
				Value: id.String(),
			}
			http.SetCookie(w, c)
			sessions[c.Value] = name

			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
	}
	tmp.ExecuteTemplate(w, "login.html", nil)
}
func home(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	if _, ok := sessions[c.Value]; !ok {
		http.Redirect(w, r, "/", http.StatusBadGateway)
	}
	tmp.ExecuteTemplate(w, "home.html", sessions[c.Value])
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	delete(sessions, cookie.Value)
	cookie = &http.Cookie{
		Name:   "Session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func myLog(message string) {
	/* ==============================writing=log==============================*/
	f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)
	log.Printf(message)
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/home", home)
	http.ListenAndServe(":5000", nil)
}
