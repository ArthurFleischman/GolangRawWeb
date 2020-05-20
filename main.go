package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/TKfleBR/GolangRawWeb/auth"
	"github.com/TKfleBR/GolangRawWeb/db"
	"github.com/TKfleBR/GolangRawWeb/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var users = map[string]models.User{}

var sessions = map[string]string{}

var tmp *template.Template

func init() {
	tmp = template.Must(template.ParseGlob("templates/*.gohtml"))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	type user struct {
		Name     string
		password []byte
	}
}

//func
func functionLogout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	delete(sessions, cookie.Value)
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func isLogged(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		log.Println("user not logged in")
		return false
	} else if _, ok := sessions[c.Value]; ok {
		return true
	} else {
		logMessage := fmt.Sprintf("invalid cookie %s\n", c.Value)
		log.Printf(logMessage)
		return false
	}
}
func myLog(message string) {
	/* ==============================writing=log==============================*/
	f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)
	log.Printf(message)
}

//views
func viewIndex(w http.ResponseWriter, r *http.Request) {
	tmp.ExecuteTemplate(w, "index.gohtml", nil)
}

func viewSignup(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		name := r.FormValue("username")
		pass := r.FormValue("password")
		cryptoPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
		u := models.User{Name: name, Password: cryptoPassword}
		err := db.InsertUser(u)
		if err != nil {
			fmt.Println("error signup")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmp.ExecuteTemplate(w, "signup.gohtml", nil)
}

func viewHome(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	if _, ok := sessions[c.Value]; !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmp.ExecuteTemplate(w, "home.gohtml", sessions[c.Value])
}
func viewLogin(w http.ResponseWriter, r *http.Request) {
	//Check if user already has a cookie(aka is logged)
	if isLogged(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	//login post request
	if r.Method == http.MethodPost {
		var u models.User
		u.Name = r.FormValue("username")
		u.Password = []byte(r.FormValue("password"))
		if !auth.User(&u, u.Password) {
			fmt.Fprintln(w, "username and/or password does not match")
			return
		}

		id := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: id.String(),
		}
		http.SetCookie(w, c)
		sessions[c.Value] = u.Name
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	}
	tmp.ExecuteTemplate(w, "login.gohtml", nil)
}

func main() {
	//getviews
	http.HandleFunc("/home", viewHome)
	http.HandleFunc("/", viewIndex)
	http.HandleFunc("/signup", viewSignup)
	http.HandleFunc("/login", viewLogin)
	// functions
	http.HandleFunc("/logout", functionLogout)
	//setup
	var port string
	if len(os.Args) < 2 {
		port = "5000"
	} else {
		port = os.Args[1]
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("given port is not valid", err)
	}
}
