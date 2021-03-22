package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is ...
type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Phone        string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	Address      string             `json:"address,omitempty" bson:"address,omitempty"`
	CurrentOrder []Product          `json:"currentorder,omitempty" bson:"currentorder,omitempty"`
	PastOrder    []Product          `json:"pastorder,omitempty" bson:"pastorder,omitempty"`
	InTransit    []Product          `json:"intransit,omitempty" bson:"intransit,omitempty"`
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
	ID1   primitive.ObjectID `json:"id" bson:"id"`
	Sub   primitive.ObjectID `json:"sub,omitempty" bson:"sub,omitempty"`
	Exist bool               `json:"exist,omitempty" bson:"exist,omitempty"`
}
type Items struct {
	Id            primitive.ObjectID   `json:"_id" bson:"_id"`
	Subcategoryid primitive.ObjectID   `json:"subcategoryid,omitempty" bson:"subcategoryid,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	Img           []string             `json:"img,omitempty" bson:"img,omitempty"`
	Details       string               `json:"details,omitempty" bson:"details,omitempty"`
	Price         int                  `json:"price,omitempty" bson:"price,omitempty"`
	Rent          int                  `json:"rent,omitempty" bson:"rent,omitempty"`
	Duration      int                  `json:"duration,omitempty" bson:"duration,omitempty"`
	Itemsid       []primitive.ObjectID `json:"itemsid,omitempty" bson:"itemsid,omitempty"`
	LocationID    primitive.ObjectID   `json:"locationid,omitempty" bson:"locationid,omitempty"`
	Stock         int                  `json:"stock,omitempty" bson:"stock,omitempty"`
	Deposit       int                  `json:"deposit,omitempty" bson:"deposit,omitempty"`
	Demand        int                  `json:"demand,omitempty" bson:"demand,omitempty"`
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
	Date     time.Time          `json:"checkoutdate,omitempty" bson:"checkoutdate,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Count    int                `json:"count,omitempty" bson:"count,omitempty"`
	Duration int                `json:"duration,omitempty" bson:"duration,omitempty"`
	Rent     int                `json:"_rent,omitempty" bson:"_rent,omitempty"`
	Deposit  int                `json:"deposit,omitempty" bson:"deposit,omitempty`
}

type Cartproduct struct {
	Status    bool               `json:"status" bson:"status"`
	Userid    primitive.ObjectID `json:"userid" bson:"userid"`
	Productid primitive.ObjectID `json:"productid" bson:"productid"`
}

//Id ...
type CartContainer struct {
	UserID primitive.ObjectID `json:"UserID" bson:"UserID"`
	ItemID primitive.ObjectID `json:"ItemID" bson:"ItemID"`
	Status bool               `json:"Status" bson:"Status"`
}

type Total struct {
	CartID primitive.ObjectID `json:"CartID" bson:"CartID"`
	UserID primitive.ObjectID `json:"UserID,omitempty" bson:"UserID,omitempty"`
	Rent   int                `json:"Rent" bson:"Rent"`
}

//SearchProduct ...
type SearchProduct struct {
	Search     string             `json:"Search,omitempty" bson:"Search,omitempty"`
	LocationID primitive.ObjectID `json:"locationid,omitempty" bson:"locationid,omitempty"`
}

type CartInput struct {
	Userid  primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Value   int                `json:"value,omitempty" bson:"value,omitempty"`
	Status  bool               `json:"status,omitempty" bson:"status,omitempty"`
	Product Product            `json:"product,omitempty" bson:"product,omitempty"`
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
