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

	//product apis for appliances
	r.HandleFunc("/api/productslist/appliances/tv", controller.ProductsListTV).Methods("POST")
	r.HandleFunc("/api/productslist/appliances/wm", controller.ProductsListWM).Methods("POST")
	r.HandleFunc("/api/productslist/appliances/fridge", controller.ProductsListFridge).Methods("POST")
	r.HandleFunc("/api/productslist/appliances/ac", controller.ProductsListAC).Methods("POST")
	r.HandleFunc("/api/productslist/appliances/others", controller.ProductsListOthApp).Methods("POST")

	//product apis for kitchen
	r.HandleFunc("/api/productslist/kitchen/stove", controller.ProductsListStove).Methods("POST")
	r.HandleFunc("/api/productslist/kitchen/micro", controller.ProductsListMicro).Methods("POST")
	r.HandleFunc("/api/productslist/kitchen/racks", controller.ProductsListRacks).Methods("POST")
	r.HandleFunc("/api/productslist/kitchen/others", controller.ProductsListOthKit).Methods("POST")

	r.HandleFunc("/api/usercreation", controller.UserCreationHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
