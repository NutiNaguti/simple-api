package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"simple-api/app"
	"simple-api/controllers"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // добавляем middleware для проверки токена
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/add-contact", controllers.AddContacts).Methods("POST")
	router.HandleFunc("/api/user/get-contact", controllers.GetContactsFor).Methods("GET")

	port := os.Getenv("PORT") // Получить порт из файла .env (при отсутствии возвращается пустая строка)
	if port == "" {
		port = "8080" // localhost
	}

	fmt.Println("Server listening")
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Println(err)
	}
}
