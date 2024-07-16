package controllers

import (
	"database/sql"
	"encoding/json"
	"gamabunta/models"
	"gamabunta/utils"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("views/login.html")
		tmpl.Execute(w, nil)
		return
	}

	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err = db.QueryRow("SELECT password FROM user WHERE username=?", creds.Username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if utils.CheckPasswordHash(creds.Password, storedPassword) {
		w.Write([]byte("Login successful"))
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("views/register.html")
		tmpl.Execute(w, nil)
		return
	}

	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(creds.Password)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Println(creds.Username)
	log.Println(creds.Password)

	_, err = db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", creds.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	w.Write([]byte("User registered successfully"))
}
