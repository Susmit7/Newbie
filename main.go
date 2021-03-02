package main

import (
	"fmt"
	"log"
	"net/http"

	controller "Newbie/controllers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	http.Handle("/", r)
	fmt.Println("Listening on Port 8080......")

	//signup and login apis
	r.HandleFunc("/api/signUp", controller.SignUpHandler).Methods("POST")
	r.HandleFunc("/api/auth", controller.AuthHandler).Methods("POST")
	r.HandleFunc("/api/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/api/resend", controller.Resendotp).Methods("GET")
	r.HandleFunc("/api/carousel", controller.Carousel).Methods("GET")
	r.HandleFunc("/api/product", controller.ProductHandler).Methods("POST")
	r.HandleFunc("/api/productslist", controller.ProductsList).Methods("POST")
	r.HandleFunc("/api/usercreation", controller.UserCreationHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
