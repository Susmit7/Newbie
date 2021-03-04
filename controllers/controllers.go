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
	authToken := ""
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

	collection, client, err := db.GetDBCollection("user")

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
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
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
			collection, client, err := db.GetDBCollection("user")
			_, err = collection.DeleteOne(context.TODO(), bson.M{"phone": userOtp.Number})
			if err != nil {
				log.Fatal(err)
			}
			trials = 0
			err = client.Disconnect(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
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
	collection, client, err := db.GetDBCollection("user")
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
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
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

//productslist APi
func ProductsList(w http.ResponseWriter, r *http.Request) {

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
	collection, client, err := db.GetDBCollection("products")

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}
	docID, err := primitive.ObjectIDFromHex(id.ID1)
	if err != nil {
		log.Fatal(err)
	}
	docID1, err := primitive.ObjectIDFromHex(id.Sub)
	if err != nil {
		log.Fatal(err)
	}
	cursor, err := collection.Find(context.TODO(), bson.M{"locationid": docID, "subcategoryid": docID1})
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Items
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var items model.Items
		if err = cursor.Decode(&items); err != nil {
			log.Fatal(err)
		}
		list = append(list, items)

	}
	json.NewEncoder(w).Encode(list)
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

//user creation

func UserCreationHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/usercreation" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "PUT" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	collection, client, err := db.GetDBCollection("user")
	if err != nil {
		log.Fatal(err)
	}

	var user model.User
	var res model.ResponseResult
	var id model.Id
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		res.Error = "Error While Creating User, Try Again"
		json.NewEncoder(w).Encode(res)
	}
	oid, _ := result.InsertedID.(primitive.ObjectID)
	id.ID1 = oid.Hex()

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	collection, client, err1 := db.GetDBCollection("wishlist")
	if err1 != nil {
		log.Fatal(err1)
	}
	var wish model.Wishlist
	wish.Userid = oid
	result1, err2 := collection.InsertOne(context.TODO(), wish)
	if err2 != nil {
		log.Fatal(err2)
	}
	oidw, _ := result1.InsertedID.(primitive.ObjectID)

	fmt.Println(oidw)

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	collection, client, err22 := db.GetDBCollection("cart")
	if err22 != nil {
		log.Fatal(err22)
	}
	result2, err33 := collection.InsertOne(context.TODO(), wish)
	if err33 != nil {
		log.Fatal(err33)
	}
	oidc, _ := result2.InsertedID.(primitive.ObjectID)

	fmt.Println(oidc)
	json.NewEncoder(w).Encode(id)
}

// wishlist api

func WishlistHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/wishlist" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//var product model.Product
	var wishlist model.Cartproduct
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &wishlist)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, client, err := db.GetDBCollection("wishlist")
	if err != nil {
		log.Fatal(err)
	}
	docID, err := primitive.ObjectIDFromHex(wishlist.Userid)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"userid": docID}
	docID1, err := primitive.ObjectIDFromHex(wishlist.Productid)
	if err != nil {
		log.Fatal(err)
	}

	if wishlist.Status == true {
		update := bson.M{"$push": bson.M{"itemsId": docID1}}
		_, err1 := collection.UpdateOne(context.TODO(), filter, update)
		if err1 != nil {
			log.Fatal(err1)
		}
		response := true
		json.NewEncoder(w).Encode(response)

	} else if wishlist.Status == false {
		update := bson.M{"$pull": bson.M{"itemsId": docID1}}
		_, err1 := collection.UpdateOne(context.TODO(), filter, update)
		if err1 != nil {
			log.Fatal(err1)
		}
		response := false
		json.NewEncoder(w).Encode(response)
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

}

//wishlist products showing api

func WishlistProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/wishlistproducts" {
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

	}
	collection, client, err := db.GetDBCollection("wishlist")

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}
	docID, err := primitive.ObjectIDFromHex(id.ID1)
	if err != nil {
		log.Fatal(err)
	}
	var product model.Wishlistarray
	err = collection.FindOne(context.TODO(), bson.M{"userid": docID}).Decode(&product)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection1, client1, err1 := db.GetDBCollection("products")
	if err1 != nil {
		res.Error = err1.Error()
		json.NewEncoder(w).Encode(res)

	}
	//fmt.Println(len(product.Wisharr))
	//fmt.Println(product.Wisharr[0])
	var list []model.Items
	var item model.Items
	for i := 0; i < len(product.Wisharr); i++ {

		err = collection1.FindOne(context.TODO(), bson.M{"_id": product.Wisharr[i]}).Decode(&item)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)

		}
		list = append(list, item)
	}
	json.NewEncoder(w).Encode(list)
	err = client1.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
