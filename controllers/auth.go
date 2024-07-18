package controllers

import (
	"encoding/json"
	"gamabunta/models"
	"gamabunta/utils"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// language=sql
	query := `
	SELECT password 
	FROM user 
	WHERE username = ?
	`

	var credentials models.Credentials

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in [ReadAll]: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		log.Printf("Error in [Unmarshal]: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var storedPassword string

	err = db.QueryRowContext(ctx, query, credentials.Email).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Credenciais inválidas.", http.StatusUnauthorized)
		return
	}

	valid := utils.CheckPasswordHash(credentials.Password, storedPassword)

	if !valid {
		http.Error(w, "Credenciais inválidas.", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// language=sql
	query := `
	INSERT INTO user (username, password)
	VALUES (?, ?)
	`

	var credentials models.Credentials

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in [ReadAll]: %v", err)
		return
	}

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		log.Printf("Error in [Unmarshal]: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validEmail := validateEmail(credentials.Email)

	err = sendVerificationEmail(credentials.Email)
	if err != nil {
		log.Printf("Error in [sendVerificationEmail]: %v", err)
		return
	}

	if !validEmail {
		http.Error(w, "E-mail inválido.", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := utils.HashPassword(credentials.Password)
	if err != nil {
		log.Printf("Error in [HashPassword]: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.ExecContext(ctx, query, credentials.Email, hashedPassword)
	if err != nil {
		log.Printf("Error in [ExecContext | query]: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("Cadastro realizado com sucesso!"))

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/register.html")
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/home.html")
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func sendVerificationEmail(email string) error {
	// Configuration
	from := "enzorocha1605@gmail.com"
	password := "emrk azlm xpvw janh"
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Obrigado por se cadastrar! Bem-vindo ao chat.")

	// Create authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
