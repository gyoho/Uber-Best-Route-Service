package models

import "gopkg.in/mgo.v2/bson"

type (
    // User represents the structure of our resource
    // need capitalized to expoert the variables
    // apply alias to use lower case names when delivering json.
    User struct {
        Id          bson.ObjectId   `json:"id" bson:"id"`
        Name        string          `json:"name" bson:"name"`
        Address     string          `json:"address" bson:"address"`
        City        string          `json:"city" bson:"city"`
        State       string          `json:"state" bson:"state"`
        Zip         string          `json:"zip" bson:"zip"`
        Coordinate  Coord           `json:"coordinate" bson:"coordinate"`
    }

    Coord struct {
        Lat float64 `json:"lat" bson:"lat"`
        Lng float64 `json:"lng" bson:"lng"`
    }
)

// append json/mgo struct tag
// to instruct how to store the user info
