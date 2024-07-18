package main

import (
	"database/sql"
	"gamabunta/controllers"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:lincetech@tcp(127.0.0.1:3306)/gamabunta")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	controllers.InitDB(db)

	// Routes
	http.HandleFunc("/", controllers.LoginPage)
	http.HandleFunc("/login", controllers.LoginPage)
	http.HandleFunc("/register", controllers.RegisterPage)
	http.HandleFunc("/home", controllers.HomePage)

	// Functions
	http.HandleFunc("/auth/login", controllers.Login)
	http.HandleFunc("/auth/register", controllers.Register)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
