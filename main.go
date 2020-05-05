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
		/* ==============================writing=log==============================*/
		f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		log.SetOutput(f)
		log.Printf("user %s registered\n", u.Name)
		/* ==============================writing=log==============================*/
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
		/* ==============================writing=log==============================*/
		f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		log.SetOutput(f)
		log.Printf("invalid cookie %s\n", c.Value)
		/* ==============================writing=log==============================*/
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
			id, err := uuid.NewV4()
			if err != nil {
				/* ==============================writing=log==============================*/
				f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				defer f.Close()
				log.SetOutput(f)
				log.Println("error", err)
				/* ==============================writing=log==============================*/
			} else {
				c := &http.Cookie{
					Name:  "session",
					Value: id.String(),
				}
				http.SetCookie(w, c)
				sessions[c.Value] = name

			}
			http.Redirect(w, r, "/home", http.StatusSeeOther)

		} else {
			/* ==============================writing=log==============================*/
			f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			defer f.Close()
			log.SetOutput(f)
			log.Printf("failed to log as %s\n", users[name].Name)
			/* ==============================writing=log==============================*/

		}
	}
	tmp.ExecuteTemplate(w, "login.html", nil)
}
func home(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	tmp.ExecuteTemplate(w, "home.html", sessions[c.Value])
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/home", home)
	http.ListenAndServe(":5000", nil)
}
