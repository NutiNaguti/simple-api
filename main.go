package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"simple-api/app"
	"simple-api/controllers"
)

func main() {
	log.SetFlags(log.Lshortfile)

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // добавляем middleware для проверки токена
	router.HandleFunc("/api/user/new", controllers.CreateAccount)
	router.HandleFunc("/api/user/login", controllers.Authenticate)

	//TODO: Доделать прототип чата
//--------------------------------------------------------------------------------------\\
//	server := services.NewServer("/entry")
//	go server.Listen()
//--------------------------------------------------------------------------------------\\

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
