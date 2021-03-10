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
	"regexp"
	"strconv"
	"strings"
	"time"

	"Newbie/db"
	model "Newbie/models"
	"Newbie/query"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	otp    string = "0000"
	trials        = 0
)

func Check(url string, method string, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/"+url {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != method {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
}

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

	Check("account", "POST", w, r)

	w.Header().Set("Content-Type", "application/json")
	var data model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}
	var user model.User

	_ = query.FindoneID("user", data.ID1, "_id").Decode(&user)
	match, err := regexp.MatchString("[0-9]{10}", user.Phone)
	fmt.Println(match)
	if data.Exist == false && user.Phone == "" {
		res.Result = "Not registered"
		json.NewEncoder(w).Encode(res)
	} else if data.Exist == false && match {
		res.Result = "Login required"
		json.NewEncoder(w).Encode(res)
	} else if data.Exist == true {
		json.NewEncoder(w).Encode(user)

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

	Check("resend", "GET", w, r)

	otpauth()
}

//carousel

func Carousel(w http.ResponseWriter, r *http.Request) {

	Check("carousel", "GET", w, r)

	var picture model.Carousel

	picture.Carousel = [7]string{"https://rht007.s3.amazonaws.com/carousel/1.jpg", "https://rht007.s3.amazonaws.com/carousel/2.jpg", "https://rht007.s3.amazonaws.com/carousel/3.jpg", "https://rht007.s3.amazonaws.com/carousel/4.jpg", "https://rht007.s3.amazonaws.com/carousel/5.jpg", "https://rht007.s3.amazonaws.com/carousel/6.jpg", "https://rht007.s3.amazonaws.com/carousel/7.jpg"}

	json.NewEncoder(w).Encode(picture)

}

//productslist APi
func ProductsList(w http.ResponseWriter, r *http.Request) {

	Check("productslist", "POST", w, r)

	w.Header().Set("Content-Type", "application/json")
	var items model.Items
	var id model.Id
	var list []model.Items
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"locationid": query.DocId(id.ID1), "subcategoryid": query.DocId(id.Sub)}

	cursor := query.FindAll("products", filter)

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		items.Img = nil
		items.Itemsid = nil
		if err = cursor.Decode(&items); err != nil {
			log.Fatal(err)
		}
		list = append(list, items)

	}
	json.NewEncoder(w).Encode(list)
}

//user creation

func UserCreationHandler(w http.ResponseWriter, r *http.Request) {

	Check("usercreation", "GET", w, r)

	var user model.User

	var id model.Id
	result := query.InsertOne("user", user)

	oid, _ := result.InsertedID.(primitive.ObjectID)
	id.ID1 = oid.Hex()

	var wish model.Wishlist
	wish.Userid = oid

	result1 := query.InsertOne("wishlist", wish)
	oidw, _ := result1.InsertedID.(primitive.ObjectID)

	fmt.Println(oidw)

	result2 := query.InsertOne("cart", wish)
	oidc, _ := result2.InsertedID.(primitive.ObjectID)

	fmt.Println(oidc)
	json.NewEncoder(w).Encode(id)
}

// wishlist api

func WishlistHandler(w http.ResponseWriter, r *http.Request) {

	Check("wishlist", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")
	var wishlist model.Cartproduct
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &wishlist)

	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"userid": query.DocId(wishlist.Userid)}

	if wishlist.Status == true {
		update := bson.M{"$push": bson.M{"itemsId": query.DocId(wishlist.Productid)}}
		query.UpdateOne("wishlist", filter, update)
		response := true
		json.NewEncoder(w).Encode(response)

	} else if wishlist.Status == false {
		update := bson.M{"$pull": bson.M{"itemsId": query.DocId(wishlist.Productid)}}
		query.UpdateOne("wishlist", filter, update)
		response := false
		json.NewEncoder(w).Encode(response)
	}
}

//wishlist products showing api

func WishlistProductsHandler(w http.ResponseWriter, r *http.Request) {

	Check("wishlistproducts", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		log.Fatal(err)

	}
	var product model.Wishlistarray
	err = query.FindoneID("wishlist", id.ID1, "userid").Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Items
	var item model.Items
	for i := 0; i < len(product.Wisharr); i++ {
		item.Img = nil
		item.Itemsid = nil
		err = query.FindoneID("products", product.Wisharr[i].Hex(), "_id").Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, item)
	}
	json.NewEncoder(w).Encode(list)
}

//product details showing api
func ProductDetailsHandler(w http.ResponseWriter, r *http.Request) {

	Check("productdetails", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var item model.Items
	err = query.FindoneID("products", id.ID1, "_id").Decode(&item)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(item)

}
