package main

import (
	controller "Newbie/controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	http.Handle("/", r)
	fmt.Println("Starting Server.....")
	fmt.Println("Listening on Port 8080......")

	r.HandleFunc("/api/account", controller.SignUpHandler).Methods("POST")
	r.HandleFunc("/api/auth", controller.AuthHandler).Methods("POST")
	r.HandleFunc("/api/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/api/resend", controller.Resendotp).Methods("GET")
	r.HandleFunc("/api/carousel", controller.Carousel).Methods("GET")

	r.HandleFunc("/api/productslist", controller.ProductsList).Methods("POST")
	r.HandleFunc("/api/usercreation", controller.UserCreationHandler).Methods("GET")
	r.HandleFunc("/api/wishlist", controller.WishlistHandler).Methods("POST")
	r.HandleFunc("/api/wishlistproducts", controller.WishlistProductsHandler).Methods("POST")

	r.HandleFunc("/api/productdetails", controller.ProductDetailsHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
