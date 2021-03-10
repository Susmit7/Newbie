package models

import (
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
	Number     string `json:"phone" bson:"phone"`
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

type Cartproduct struct {
	Status    bool   `json:"status" bson:"status"`
	Userid    string `json:"userid" bson:"userid"`
	Productid string `json:"productid" bson:"productid"`
}

type Wishlistarray struct {
	Wisharr []primitive.ObjectID `json:"itemsid" bson:"itemsid"`
}

type Account struct {
	ID    string `json:"userid"`
	Exist bool   `json:"exist"`
}
