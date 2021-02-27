package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"Newbie/db"
	model "Newbie/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	otp    string = "0000"
	trials        = 0
)

func otpauth() {
	accountSid := "AC1cab9315c49a09f2e53bea328a4799bf"
	authToken := "3454674edbd72c9d656a479c80495ad0"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/AC1cab9315c49a09f2e53bea328a4799bf/Messages.json"

	max := 9999
	min := 1000
	rand.Seed(time.Now().UnixNano())
	otp = strconv.Itoa(rand.Intn(max-min+1) + min)

	msgData := url.Values{}
	msgData.Set("To", "+918001044568")
	msgData.Set("From", "+18174426921")
	msgData.Set("Body", otp)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)

		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}

// SignUpHandler ...
func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/signUp" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	collection, err := db.GetDBCollection("user")

	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{Key: "phone", Value: user.Phone}}).Decode(&result)

	if err == nil {
		res.Result = "Phone Number already Registered!!"
		json.NewEncoder(w).Encode(res)
		return
	}

	if err != nil {
		if err.Error() == "mongo: no documents in result" {

			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}

			otpauth()

			res.Result = "Phone Authentication Required!"
			json.NewEncoder(w).Encode(res)
			return
		}
	}

}

//SignUpAuthHandler ...
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/auth" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var userOtp model.OtpContainer
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &userOtp)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if userOtp.OtpEntered == otp && userOtp.From == "signup" {
		fmt.Println("The signUp authentication is successful!")
		res.Result = "The signUp authentication is successful!"
		json.NewEncoder(w).Encode(res)

	} else if userOtp.OtpEntered != otp && userOtp.From == "login" {
		res.Error = "OTP Did not Match!"
		json.NewEncoder(w).Encode(res)

	} else if userOtp.OtpEntered == otp && userOtp.From == "login" {
		fmt.Println("The Login authentication is successful!")
		res.Result = "The Login authentication is successful!"
		json.NewEncoder(w).Encode(res)

	} else if userOtp.OtpEntered != otp && userOtp.From == "signup" {
		if trials < 5 {
			res.Error = "OTP Did not Match!"
			json.NewEncoder(w).Encode(res)
			trials++

		} else if trials == 5 {
			res.Error = "Data Deleted, Signup again!"
			json.NewEncoder(w).Encode(res)
			collection, err := db.GetDBCollection("user")
			_, err = collection.DeleteOne(context.TODO(), bson.M{"phone": userOtp.Number})
			if err != nil {
				log.Fatal(err)
			}
			trials = 0
		}

	}

}

//login handler

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var login model.Login
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &login)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, err := db.GetDBCollection("user")
	var result model.Login
	err = collection.FindOne(context.TODO(), bson.D{{Key: "phone", Value: login.Contact}}).Decode(&result)

	if err != nil {
		res.Result = "Not Registered!"
		json.NewEncoder(w).Encode(res)
		return
	}

	if err == nil {
		res.Result = "Welcome Buddy,Enter Otp!"
		json.NewEncoder(w).Encode(res)

		otpauth()

	}
}

//resend otp

func Resendotp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/resend" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	otpauth()
}

//carousel

func Carousel(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/carousel" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var picture model.Carousel

	picture.Carousel = [7]string{"https://rht007.s3.amazonaws.com/carousel/1.jpg", "https://rht007.s3.amazonaws.com/carousel/2.jpg", "https://rht007.s3.amazonaws.com/carousel/3.jpg", "https://rht007.s3.amazonaws.com/carousel/4.jpg", "https://rht007.s3.amazonaws.com/carousel/5.jpg", "https://rht007.s3.amazonaws.com/carousel/6.jpg", "https://rht007.s3.amazonaws.com/carousel/7.jpg"}

	json.NewEncoder(w).Encode(picture)

}

//product_details handler

func ProductHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/product" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var product model.Product
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	collection, err := db.GetDBCollection("products")

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	docID, err := primitive.ObjectIDFromHex(id.ID)

	err = collection.FindOne(context.TODO(), bson.M{"_id": docID}).Decode(&product)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	json.NewEncoder(w).Encode(product)

}

//productlisthandler
func ProductsListHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/productslist" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//var product model.Product
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, err := db.GetDBCollection("products")

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}
	docID, err := primitive.ObjectIDFromHex(id.ID)
	cursor, err := collection.Find(context.TODO(), bson.M{"locationid": docID})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var items model.Items
		if err = cursor.Decode(&items); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(items)

	}
}
