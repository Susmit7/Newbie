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
	"go.mongodb.org/mongo-driver/mongo/options"
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
	accountSid := "######"
	authToken := "#####"
	urlStr := "######"

	max := 9999
	min := 1000
	rand.Seed(time.Now().UnixNano())
	otp = strconv.Itoa(rand.Intn(max-min+1) + min)

	msgData := url.Values{}
	msgData.Set("To", "####")
	msgData.Set("From", "######")
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
	var data model.Account
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}
	var user model.User

	_ = query.FindoneID("user", data.ID.Hex(), "_id").Decode(&user)
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
	Check("auth", "POST", w, r)
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

	Check("login", "POST", w, r)

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
	collection, client := query.Connection("products")
	for i := 0; i < len(product.Wisharr); i++ {
		item.Img = nil
		item.Itemsid = nil
		//_ = query.FindoneID("products", product.Wisharr[i].Hex(), "_id").Decode(&item)
		_ = collection.FindOne(context.TODO(), bson.M{"_id": product.Wisharr[i]}).Decode(&item)
		list = append(list, item)
	}
	defer query.Endconn(client)
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

//checkout api

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	Check("checkout", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	filter := bson.M{"userid": query.DocId(id.ID1)}
	update := bson.M{"$set": bson.M{"product": bson.A{}}}

	var products model.Cart

	_ = query.FindoneID("cart", id.ID1, "userid").Decode(&products)

	query.UpdateOne("cart", filter, update)

	var count []int
	var productid []primitive.ObjectID
	count = nil
	productid = nil
	for i := 0; i < len(products.Product); i++ {
		count = append(count, products.Product[i].Count)
		productid = append(productid, products.Product[i].P_id)
	}

	filter1 := bson.M{"_id": query.DocId(id.ID1)}
	collection, client := query.Connection("user")
	for i := 0; i < len(products.Product); i++ {
		update1 := bson.M{"$push": bson.M{"currentorder": products.Product[i]}}
		_, err := collection.UpdateOne(context.TODO(), filter1, update1)
		if err != nil {
			log.Fatal(err)
		}
	}
	query.Endconn(client)

	fmt.Println(productid)
	fmt.Println(count)
	collection2, client2 := query.Connection("user")
	for i := 0; i < len(products.Product); i++ {
		filter2 := bson.M{"_id": productid[i]}
		update2 := bson.M{"$inc": bson.M{"stock": -count[i]}}
		_, err := collection2.UpdateOne(context.TODO(), filter2, update2, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatal(err)
		}
	}
	query.Endconn(client2)

}

//update cart api
func UpdateCart(w http.ResponseWriter, r *http.Request) {

	Check("updatecart", "PUT", w, r)

	w.Header().Set("Content-Type", "application/json")
	var cart model.CartContainer
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &cart)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	filter := bson.M{"userid": query.DocId(cart.UserID)}

	if bool(cart.Status) == true {
		update := bson.M{"$push": bson.M{"itemsId": query.DocId(cart.ItemID)}}
		query.UpdateOne("cart", filter, update)

		response := true
		json.NewEncoder(w).Encode(response)

	} else if bool(cart.Status) == false {
		update1 := bson.M{"$push": bson.M{"itemsId": query.DocId(cart.ItemID)}}

		query.UpdateOne("cart", filter, update1)

		response := false
		json.NewEncoder(w).Encode(response)
	}
}

//SearcEngine api
func SearchEngine(w http.ResponseWriter, r *http.Request) {

	Check("searchengine", "POST", w, r)

	w.Header().Set("Content-Type", "application/json")

	body, _ := ioutil.ReadAll(r.Body)

	var srch model.SearchProduct

	var res model.ResponseResult

	err := json.Unmarshal(body, &srch)
	if err != nil {
		log.Fatal(w, "error occured while unmarshling")
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	}

	search := bson.M{"$text": bson.M{"$search": srch.Search}}

	cursor := query.FindAll("products", search)

	var show []model.Items
	var product model.Items
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		product.Img = nil
		product.Itemsid = nil
		if err = cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		show = append(show, product)
	}

	json.NewEncoder(w).Encode(show)

}

//cart products api
func CartProducts(w http.ResponseWriter, r *http.Request) {

	Check("cartproducts", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	}

	var item model.Cart

	err = query.FindoneID("cart", id.ID1, "userid").Decode(&item)
	if err != nil {
		log.Fatal(err)
	}
	var prod []model.Product
	for i := 0; i < len(item.Product); i++ {
		prod = append(prod, item.Product[i])
	}

	json.NewEncoder(w).Encode(prod)

}

//Cart First Time
func CartFirstTime(w http.ResponseWriter, r *http.Request) {

	Check("cartfirsttime", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var ct model.CartInput
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ct)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, client := query.Connection("cart")

	var doc model.Cart

	err = collection.FindOne(context.TODO(), bson.M{"userid": ct.Userid, "product.p_id": ct.Product.P_id, "product.duration": ct.Product.Duration}).Decode(&doc)

	if err != nil {

		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid}, bson.M{"$push": bson.M{"product": ct.Product}})

		res1 := "New product added"
		json.NewEncoder(w).Encode(res1)

	} else {

		update := bson.M{"$set": bson.M{"product.$.count": ct.Product.Count + 1}}

		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid, "product.p_id": ct.Product.P_id, "product.duration": ct.Product.Duration}, update)
		if err != nil {
			log.Fatal(err)
		}
		res2 := "Count of product increased"
		json.NewEncoder(w).Encode(res2)

	}
	query.Endconn(client)

}

//CartInput

func CartInput(w http.ResponseWriter, r *http.Request) {

	Check("cartinput", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var ct model.CartInput
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ct)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	}
	collection, client := query.Connection("cart")

	var doc model.Cart

	err = collection.FindOne(context.TODO(), bson.M{"userid": ct.Userid, "product.p_id": ct.Product.P_id, "product.count": ct.Product.Count, "product.duration": ct.Product.Duration}).Decode(&doc)
	if err != nil {
		//didn't found any match
		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid}, bson.M{"$push": bson.M{"product": ct.Product}})

		respn := "New Product Added"
		json.NewEncoder(w).Encode(respn)

	} else {
		// if found the match
		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid, "product.p_id": ct.Product.P_id}, bson.M{"$set": bson.M{"product.count": ct.Product.Count, "product.duration": ct.Product.Duration, "product._rent": ct.Product.Rent}})

		respm := "Existing Product Updated"
		json.NewEncoder(w).Encode(respm)

	}
	query.Endconn(client)
}

//Remove Cart Products

func RemoveCartProduct(w http.ResponseWriter, r *http.Request) {

	Check("removecartproduct", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var rem model.RemoveCartProduct

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &rem)

	filter := bson.M{"userid": rem.UserId}
	update := bson.M{"$pull": bson.M{"product": bson.M{"p_id": rem.ProductId, "count": rem.Count, "duration": rem.Duration}}}
	query.UpdateOne("cart", filter, update)

	if err != nil {
		log.Fatal(err)
	}

	respn := "Data Removed"
	json.NewEncoder(w).Encode(respn)
}

//to be changed
func ProductStock(w http.ResponseWriter, r *http.Request) {

	Check("stock", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var s model.StockId
	var sd model.StockData

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &s)

	var res model.ResponseResult

	err = query.FindoneID("products", s.ProductId.Hex(), "_id").Decode(&sd)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}

	json.NewEncoder(w).Encode(sd.Stock)
}
