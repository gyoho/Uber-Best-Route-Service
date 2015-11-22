package models

import "gopkg.in/mgo.v2/bson"

type Trip struct {
    Id                              bson.ObjectId   `json:"id" bson:"id"`
    Status                          string          `json:"status" bson:"status"`
    Starting_from_location_id       string          `json:"starting_from_location_id" bson:"starting_from_location_id"`
    Next_destination_location_id    string          `json:"next_destination_location_id,omitempty" bson:"next_destination_location_id"`
    Counter                         int             `json:"-" bson:"counter"`
    Best_route_location_ids         []string        `json:"best_route_location_ids" bson:"best_route_location_ids"`
    Product_ids                     []string        `json:"-" bson:"product_ids"`
    Total_uber_costs                float64         `json:"total_uber_costs" bson:"total_uber_costs"`
    Total_uber_duration             float64         `json:"total_uber_duration" bson:"total_uber_duration"`
    Total_distance                  float64         `json:"total_distance" bson:"total_distance"`
    Uber_wait_time_eta              float64         `json:"uber_wait_time_eta,omitempty" bson:"uber_wait_time_eta"`
}

// use omitempty to not display empty value in response
// **Caution: Do Not Add Space before "omitempty" tag!!**

type RideRequest struct {
    ProductID                       string           `json:"product_id"`
    StartLat                        float64          `json:"start_latitude"`
    StartLng                        float64          `json:"start_longitude"`
    EndLat                          float64          `json:"end_latitude"`
    EndLng                          float64          `json:"end_longitude"`
}
