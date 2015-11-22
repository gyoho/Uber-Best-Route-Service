package controllers

import (
    // Standard library packages
    "fmt"
    "net/http"
    "encoding/json"
    "errors"
    "log"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    // Third party packages
    "github.com/julienschmidt/httprouter"
    // Use defined model in another dir
    "../models"
    "../utils"
)

// UserController represents the controller for operating on the User resource
// i.e) controller has no property, only methods
type TripController struct{
    // use reference to access mongodb
    session *mgo.Session
}

// Constructor
func NewTripController(s *mgo.Session) *TripController {
    // instantiate with the session received as an arg
    return &TripController{s}
}

func (tc TripController) CreateTrip(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    trip, err := createTrip(tc, req)

    // Create response
    // Write content-type, statuscode, payload
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        // Marshal provided interface into JSON structure
        tripJson, _ := json.Marshal(trip)
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", tripJson)
    }
}

func createTrip(tc TripController, req *http.Request) (models.Trip, error) {
    // Stub an trip to be populated from the body
    location := utils.Location{}
    // Populate the user data
    err := json.NewDecoder(req.Body).Decode(&location)
    if err != nil {
        log.Println(err)
        return models.Trip{}, err
    }

    // Get the trip plan
    trip, err := planTrip("breadthFirst", tc, location)
    if err != nil {
        log.Println(err)
        return models.Trip{}, err
    }

    // Persist the data to mongodb
    conn := tc.session.DB("cmpe273_asgmt2").C("trip")
    err = conn.Insert(trip)
    if err != nil {
        log.Println(err)
        return models.Trip{}, err
    }

    return trip, nil
}

func (tc TripController) GetTripDetails(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    trip, err := retriveTripDetailsById(tc, param.ByName("id"))

    if err != nil {
        log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        // Create response
        // Marshal provided interface into JSON structure
        tripJson, _ := json.Marshal(trip)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", tripJson)
    }
}

func (tc TripController) UpdateTripStatus(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    updatedTripStatus, err := goToNextDestination(tc, param.ByName("id"))
    if err != nil {
        log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        // Create response
        // Marshal provided interface into JSON structure
        tripJson, _ := json.Marshal(updatedTripStatus)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", tripJson)
    }
}

func planTrip(algo string, tc TripController, location utils.Location) (models.Trip, error) {
    switch algo {
    case "breadthFirst":
        return BreadthFirst(tc, location)
    case "depthFirst":
        return DepthFirst(tc, location)
    default:
        return BreadthFirst(tc, location)
    }
}

// ******** Trip Helper Function

func BreadthFirst(tc TripController, location utils.Location) (models.Trip, error) {
    // Get all possible destination permutation
    allPossibleRoutes := utils.FindAllPossibleRoute(location.Location_ids)
    // Append start location to the first and last
    for i, _ := range allPossibleRoutes {
        allPossibleRoutes[i].Ids = append([]string{location.Starting_from_location_id}, allPossibleRoutes[i].Ids...)
        allPossibleRoutes[i].Ids = append(allPossibleRoutes[i].Ids, location.Starting_from_location_id)
    }

    // Get coordinate for each location
    var allPossibleEstimates utils.ListEstimates
    for _, route := range allPossibleRoutes {
        estimatesArray, err := formLocationArray(tc, route.Ids)
        if err != nil {
            return models.Trip{}, err
        }
        allPossibleEstimates = append(allPossibleEstimates, estimatesArray)
    }

    // Get Estimate for all routes
    trip, _ := utils.FindBestRouteBreadthFirst(allPossibleEstimates)

    // remove the home id from the best_route_location_ids (last element)
    trip.Best_route_location_ids = trip.Best_route_location_ids[0:len(trip.Best_route_location_ids)-1]

    // assign already determined values
    trip.Id = bson.NewObjectId()
    trip.Status = "planning"
    trip.Starting_from_location_id = location.Starting_from_location_id

    return trip, nil
}

func DepthFirst(tc TripController, location utils.Location) (models.Trip, error) {
    location_ids := append([]string{location.Starting_from_location_id}, location.Location_ids...)

    // Get coordinate for each location
    estimatesArray, err := formLocationArray(tc, location_ids)
    if err != nil {
        return models.Trip{}, err
    }

    // Save the starting point. Need
    startLoc := estimatesArray[0]

    // Going to update Trip plan
    trip := models.Trip{}
    trip, estimatesArray, err = utils.FindBestRouteDepthFirst(estimatesArray)
    if err != nil {
        return models.Trip{}, err
    }

    // Back to home: now the array only has one element
    estimatesArray = append(estimatesArray, startLoc)
    lastLoc := estimatesArray[0]
    estimatesArray = estimatesArray[1:]
    _, err = utils.FindRouteDepthFirst(lastLoc, estimatesArray, &trip)
    if err != nil {
        return models.Trip{}, err
    }

    // remove the home id from the best_route_location_ids
    trip.Best_route_location_ids = trip.Best_route_location_ids[:len(trip.Best_route_location_ids)-1]

    // assign already determined values
    trip.Id = bson.NewObjectId()
    trip.Status = "planning"
    trip.Starting_from_location_id = location.Starting_from_location_id

    return trip, nil
}

func formLocationArray(tc TripController, location_ids []string) (utils.Estimates, error) {
    if len(location_ids) < 1 {
        return utils.Estimates{}, errors.New("Need at least one Destination Loation")
    }

    var estimatesArray utils.Estimates
    for _, location_id := range location_ids {
        estimate := utils.Estimate{}
        estimate.Location_id = location_id

        coord, err := retriveCoordById(tc, estimate.Location_id)
        if err != nil {
            return utils.Estimates{}, err
        }
        estimate.Coordinate = coord
        estimatesArray = append(estimatesArray, estimate)
    }

    return estimatesArray, nil
}

// ********

func retriveCoordById(tc TripController, id string) (models.Coord, error) {
    // Verify id is ObjectId, otherwise bail
    if !bson.IsObjectIdHex(id) {
        return models.Coord{}, errors.New("Not valid Location ID")
    }

    // Grab id
    objId := bson.ObjectIdHex(id)
    // Stub user
    usr := models.User{}
    // make connection
    conn := tc.session.DB("cmpe273_asgmt2").C("user")
    // the id is created by system side not db side
    err := conn.Find(bson.M{"id": objId}).One(&usr)
    if err != nil {
        return models.Coord{}, errors.New("No Location found with this ID")
    }

    return usr.Coordinate, nil
}

func retriveTripDetailsById(tc TripController, id string) (models.Trip, error) {
    // Verify id is ObjectId, otherwise bail
    if !bson.IsObjectIdHex(id) {
        return models.Trip{}, errors.New("Not valid trip ID")
    }

    // Grab id
    objId := bson.ObjectIdHex(id)
    // Stub user
    trip := models.Trip{}
    // make connection
    conn := tc.session.DB("cmpe273_asgmt2").C("trip")
    // the id is created by system side not db side
    err := conn.Find(bson.M{"id": objId}).One(&trip)
    if err != nil {
        return models.Trip{}, errors.New("No trip found with this ID")
    }

    return trip, nil
}

func goToNextDestination(tc TripController, id string) (models.Trip, error) {
    // Get the trip plan
    updatedTrip, err := retriveTripDetailsById(tc, id)
    if err != nil {
        return models.Trip{}, err
    }

    if updatedTrip.Counter < len(updatedTrip.Best_route_location_ids) {
        // change the status
        if updatedTrip.Counter == 0 {
            updatedTrip.Status = "requesting"
        }

        // Update the next_destination_location_id
        updatedTrip.Next_destination_location_id = updatedTrip.Best_route_location_ids[updatedTrip.Counter]

        // Update the uber_wait_time_eta
        eta, err := requestRide(tc, updatedTrip)
        if err != nil {
            log.Println(err)
            return models.Trip{}, err
        }
        updatedTrip.Uber_wait_time_eta = eta

        // Increment the counter
        updatedTrip.Counter++

        // Grab id
        objId := bson.ObjectIdHex(id)
        // make connection
        conn := tc.session.DB("cmpe273_asgmt2").C("trip")
        err = conn.Update(bson.M{"id": objId}, updatedTrip)
        if err != nil {
            log.Println(err)
            return models.Trip{}, err
        }

        return updatedTrip, nil
    } else if updatedTrip.Counter == len(updatedTrip.Best_route_location_ids) {
        // Going back to home
        updatedTrip.Status = "finishing"
        updatedTrip.Next_destination_location_id = updatedTrip.Starting_from_location_id

        // Clean up
        objId := bson.ObjectIdHex(id)
        conn := tc.session.DB("cmpe273_asgmt2").C("trip")
        err = conn.Remove(bson.M{"id": objId})
        if err != nil {
            log.Println(err)
            return models.Trip{}, err
        }

        return updatedTrip, nil
    } else {
        return models.Trip{}, errors.New("Bug in code")
    }
}

func requestRide(tc TripController, trip models.Trip) (float64, error)  {
    reqBody := models.RideRequest{}
    var startCoord models.Coord
    var err error

    if trip.Counter == 0 {
        startCoord, err = retriveCoordById(tc, trip.Starting_from_location_id)
        if err != nil {
            return 0.0, err
        }
    } else {
        startCoord, err = retriveCoordById(tc, trip.Best_route_location_ids[trip.Counter-1])
        if err != nil {
            return 0.0, err
        }
    }

    endCoord, err := retriveCoordById(tc, trip.Next_destination_location_id)
    if err != nil {
        return 0.0, err
    }

    reqBody.ProductID = trip.Product_ids[trip.Counter]
    reqBody.StartLat = startCoord.Lat
    reqBody.StartLng = startCoord.Lng
    reqBody.EndLat = endCoord.Lat
    reqBody.EndLng = endCoord.Lng

    return utils.RequestRide(reqBody)
}
