package controllers

import (
    // Standard library packages
    "fmt"
    "net/http"
    "encoding/json"
    "math/rand"
    "errors"
    "log"
    "strconv"
    "io"

    // Third party packages
    "github.com/julienschmidt/httprouter"
    // Use defined model in another dir
    "../models"
)

// UserController represents the controller for operating on the User resource
// i.e) controller has no property, only methods
type UserController struct{
}

// returns the reference so to use its methods
func NewUserController() *UserController {
    return &UserController{}
}


// userID -> userObj
var UserMap map[int]models.User = make(map[int]models.User)

func (uc UserController) CreateUser(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    // Stub an user to be populated from the body
    usr := models.User{}
    // Populate the user data
    json.NewDecoder(req.Body).Decode(&usr)
    // assign unique ID
    usr.Id = rand.Int()

    // Append coordinate to the user object
    err := getCoordinates(&usr)
    if err != nil {
		log.Println(err)
	}


    // TODO: Persist the data
    UserMap[usr.Id] = usr


    // Create response
    // Marshal provided interface into JSON structure
    usrJson, _ := json.Marshal(usr)
    // Write content-type, statuscode, payload
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", usrJson)
    }
}

func (uc UserController) GetUser(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    usr, err := retriveUserById(param.ByName("id"))
    if err != nil {
		log.Println(err)
	}

    // Create response
    // Marshal provided interface into JSON structure
    usrJson, _ := json.Marshal(usr)
    // Write content-type, statuscode, payload
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", usrJson)
    }
}

func (uc UserController) UpdateUser(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
    updatedUsr, err := updateUserLocation(param.ByName("id"), req.Body)
    if err != nil {
		log.Println(err)
	}

    // update the hashmap
    UserMap[updatedUsr.Id] = updatedUsr

    // Create response
    // Marshal provided interface into JSON structure
    usrJson, _ := json.Marshal(updatedUsr)
    // Write content-type, statuscode, payload
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", usrJson)
    }
}

func (uc UserController) RemoveUser(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    usr, err := retriveUserById(param.ByName("id"))
    if err != nil {
		log.Println(err)
	}

    // remove the user from the hashmap
    delete(UserMap, usr.Id)

    rw.Header().Set("Content-Type", "plain/text")

    if err != nil {
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "User with ID = , %s, Deleted", param.ByName("id"))
    }
}

func retriveUserById(idStr string) (models.User, error) {
    // find the user using hashmap
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return models.User{}, errors.New("ID must be an ineger")
    }

    // TODO: access DB

    // check if the ID exists in the hashmap
    usr, ok := UserMap[id]
    if !ok {
        return models.User{}, errors.New("No user found with this ID")
    }

    return usr, nil
}

func updateUserLocation(idStr string, contents io.Reader) (models.User, error) {
    usr, err := retriveUserById(idStr)
    if err != nil {
        return models.User{}, err
    }

    updatedUsr := models.User{}
    json.NewDecoder(contents).Decode(&updatedUsr)
    // preseve ID and Name
    updatedUsr.Id = usr.Id
    updatedUsr.Name = usr.Name

    // Append coordinate to the user object
    err = getCoordinates(&updatedUsr)
    if err != nil {
        return models.User{}, err
    }

    return updatedUsr, nil
}
