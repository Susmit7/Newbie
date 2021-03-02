package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is ...
type User struct {
	Name          string             `json:"name" bson:"name,omitempty"`
	Phone         string             `json:"phone" bson:"phone,omitempty"`
	Email         string             `json:"email" bson:"email,omitempty"`
	Address       string             `json:"address" bson:"address,omitempty"`
	LocationID    primitive.ObjectID `json:"location" bson:"location,omitempty"`
	Current_order []string           `json:"current" bson:"current,omitempty"`
	Past_order    []string           `json:"past" bson:"past,omitempty"`
}

//login...

type Login struct {
	Contact string `json:"contact"`
}

// ResponseResult is ...
type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

// OtpContainer ...
type OtpContainer struct {
	OtpEntered string `json:"otpentered"`
	Number     string `json:"number,omitempty"`
	From       string `json:"from"`
}

type Carousel struct {
	Carousel [7]string `json:"carousel"`
}

type Datalist struct {
	Alldata []string `json:"alldata"`
}

type Id struct {
	ID string `json:"id" bson:"id"`
}
type Product struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name"`
	Details  string             `json:"details"`
	Rent     int                `json:"rent"`
	Duration int                `json:"duration"`
	IMG      string             `json:"img"`
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
