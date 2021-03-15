package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is ...
type User struct {
	ID           primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string               `json:"name" bson:"name,omitempty"`
	Phone        string               `json:"phone" bson:"phone,omitempty"`
	Email        string               `json:"email" bson:"email,omitempty"`
	Address      string               `json:"address" bson:"address,omitempty"`
	CurrentOrder []primitive.ObjectID `json:"currentorder" bson:"currentorder,omitempty"`
	PastOrder    []primitive.ObjectID `json:"pastorder" bson:"pastorder,omitempty"`
	InTransit    []primitive.ObjectID `json:"intransit,omitempty" bson:"intransit,omitempty"`
}

//login...

type Login struct {
	Contact string `json:"contact"`
}

// ResponseResult is ...
type ResponseResult struct {
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

// OtpContainer ...
type OtpContainer struct {
	OtpEntered string `json:"otpentered,omitempty" bson:"otp,omitempty"`
	Number     string `json:"phone,omitempty" bson:"phone,omitempty"`
	From       string `json:"from,omitempty" bson:"from,omitempty"`
}

type Carousel struct {
	Carousel [7]string `json:"carousel"`
}

type Id struct {
	ID1   string `json:"id" bson:"id"`
	Sub   string `json:"sub,omitempty" bson:"sub,omitempty"`
	Exist bool   `json:"exist,omitempty" bson:"exist,omitempty"`
}
type Items struct {
	Id            primitive.ObjectID   `json:"_id" bson:"_id"`
	Subcategoryid primitive.ObjectID   `json:"subcategoryid" bson:"subcategoryid"`
	Name          string               `json:"name" bson:"name"`
	Img           []string             `json:"img" bson:"img"`
	Details       string               `json:"details" bson:"details"`
	Price         int                  `json:"price" bson:"price"`
	Rent          int                  `json:"rent" bson:"rent"`
	Duration      int                  `json:"duration" bson:"duration"`
	Itemsid       []primitive.ObjectID `json:"itemsid" bson:"itemsid"`
	LocationID    primitive.ObjectID   `json:"locationid" bson:"locationid"`
}

type Wishlist struct {
	Userid primitive.ObjectID `json:"userid" bson:"userid"`
}

type Wishlistarray struct {
	Wisharr []primitive.ObjectID `json:"itemsid" bson:"itemsid"`
}

type Account struct {
	ID    primitive.ObjectID `json:"id"`
	Exist bool               `json:"exist"`
}

type Cart struct {
	Userid  primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Product []Product          `json:"product,omitempty" bson:"product,omitempty"`
}

type Product struct {
	P_id     primitive.ObjectID `json:"p_id,omitempty" bson:"p_id,omitempty"`
	Location primitive.ObjectID `json:"locationid,omitempty" bson:"locationid,omitempty"`
	Img      string             `json:"img,omitempty" bson:"img,omitempty"`
	Date     time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Count    int                `json:"count,omitempty" bson:"count,omitempty"`
	Duration int                `json:"duration,omitempty" bson:"duration,omitempty"`
	Rent     int                `json:"_rent,omitempty" bson:"_rent,omitempty"`
}

type Cartproduct struct {
	Status    bool   `json:"status" bson:"status"`
	Userid    string `json:"userid" bson:"userid"`
	Productid string `json:"productid" bson:"productid"`
}

//Id ...
type CartContainer struct {
	UserID string `json:"UserID" bson:"UserID"`
	ItemID string `json:"ItemID" bson:"ItemID"`
	Status bool   `json:"Status" bson:"Status"`
}

type Total struct {
	CartID string `json:"CartID" bson:"CartID"`
	UserID string `json:"UserID,omitempty" bson:"UserID,omitempty"`
	Rent   int    `json:"Rent" bson:"Rent"`
}

//SearchProduct ...
type SearchProduct struct {
	Search     string `json:"Search,omitempty" bson:"Search,omitempty"`
	LocationID string `json:"locationid,omitempty" bson:"locationid,omitempty"`
}

type CartInput struct {
	Userid  primitive.ObjectID `json:"userid" bson:"userid"`
	Product Product            `json:"product" bson:"product"`
}

type RemoveCartProduct struct {
	UserId    primitive.ObjectID `json:"userid" bson:"userid"`
	ProductId primitive.ObjectID `json:"p_id" bson:"p_id"`

	Count    int `json:"count" bson:"count"`
	Duration int `json:"duration" bson:"duration"`
	Rent     int `json:"_rent,omitempty" bson:"_rent",omitempty`
}

//ProductStock

type StockId struct {
	ProductId primitive.ObjectID `json:"_id" bson:"_id"`
}

type StockData struct {
	Id            primitive.ObjectID   `json:"_id" bson:"_id"`
	Subcategoryid primitive.ObjectID   `json:"subcategoryid" bson:"subcategoryid"`
	Name          string               `json:"name" bson:"name"`
	Img           []string             `json:"img" bson:"img"`
	Details       string               `json:"details" bson:"details"`
	Price         int                  `json:"price" bson:"price"`
	Rent          int                  `json:"rent" bson:"rent"`
	Duration      int                  `json:"duration" bson:"duration"`
	Itemsid       []primitive.ObjectID `json:"itemsid,omitempty" bson:"itemsid,omitempty"`
	LocationID    primitive.ObjectID   `json:"locationid" bson:"locationid"`
	Stock         int                  `json:"stock" bson:"stock"`
}
