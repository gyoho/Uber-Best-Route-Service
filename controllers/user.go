package controllers

import (
    // Standard library packages
    "fmt"
    "net/http"
    "encoding/json"
    "errors"
    "log"
    "io"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    // Third party packages
    "github.com/julienschmidt/httprouter"
    // Use defined model in another dir
    "../models"
)

// UserController represents the controller for operating on the User resource
// i.e) controller has no property, only methods
type UserController struct{
    // use reference to access mongodb
    session *mgo.Session
}

// Constructor
func NewUserController(s *mgo.Session) *UserController {
    // instantiate with the session received as an arg
    return &UserController{s}
}

func (uc UserController) CreateUser(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    usr, err := createUser(uc, req)

    // Create response
    // Write content-type, statuscode, payload
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        // Marshal provided interface into JSON structure
        usrJson, _ := json.Marshal(usr)
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", usrJson)
    }
}

func createUser(uc UserController, req *http.Request) (models.User, error){
    // Stub an user to be populated from the body
    usr := models.User{}

    // Populate the user data
    err := json.NewDecoder(req.Body).Decode(&usr)
    if err != nil {
        log.Println(err)
        return models.User{}, err
    }

    // assign unique string ID
    usr.Id = bson.NewObjectId()

    // Append coordinate to the user object
    err = getCoordinates(&usr)
    if err != nil {
		log.Println(err)
        return models.User{}, err
	}

    // Persist the data to mongodb
    conn := uc.session.DB("cmpe273_asgmt2").C("user")
    err = conn.Insert(usr)
    if err != nil {
        log.Println(err)
        return models.User{}, err
    }

    return usr, nil
}

func (uc UserController) GetUser(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    usr, err := retriveUserById(uc, param.ByName("id"))
    if err != nil {
		log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
	} else {
        // Create response
        // Marshal provided interface into JSON structure
        usrJson, _ := json.Marshal(usr)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", usrJson)
    }
}

func (uc UserController) UpdateUser(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
    updatedUsr, err := updateUserLocation(uc, param.ByName("id"), req.Body)
    if err != nil {
		log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
	} else {
        // Create response
        // Marshal provided interface into JSON structure
        usrJson, _ := json.Marshal(updatedUsr)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", usrJson)
    }
}

func (uc UserController) RemoveUser(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    rw.Header().Set("Content-Type", "plain/text")

    // check if the user exists in db
    _, err := retriveUserById(uc, param.ByName("id"))
    if err != nil {
		log.Println(err)
        // Write content-type, statuscode, payload
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
	} else {
        // remove the user from the list
        objId := bson.ObjectIdHex(param.ByName("id"))
        conn := uc.session.DB("cmpe273_asgmt2").C("user")
        err = conn.Remove(bson.M{"id": objId})
        if err != nil {
            log.Println(err)
            rw.WriteHeader(400)
            fmt.Fprintf(rw, "%s\n", err)
        } else {
            rw.WriteHeader(200)
            fmt.Fprintf(rw, "User with ID=%s Deleted", param.ByName("id"))
        }
    }
}

func retriveUserById(uc UserController, id string) (models.User, error) {
    // Verify id is ObjectId, otherwise bail
    if !bson.IsObjectIdHex(id) {
        return models.User{}, errors.New("Not valid user ID")
    }

    // Grab id
    objId := bson.ObjectIdHex(id)
    // Stub user
    usr := models.User{}
    // make connection
    conn := uc.session.DB("cmpe273_asgmt2").C("user")
    // the id is created by system side not db side
    err := conn.Find(bson.M{"id": objId}).One(&usr)
    if err != nil {
        return models.User{}, errors.New("No user found with this ID")
    }

    return usr, nil
}

func updateUserLocation(uc UserController, id string, contents io.Reader) (models.User, error) {
    // check if the user exists in db
    usr, err := retriveUserById(uc, id)
    if err != nil {
        return models.User{}, err
    }

    // get the updated contents
    updatedUsr := models.User{}
    updatedUsr.Id = usr.Id
    updatedUsr.Name = usr.Name
    json.NewDecoder(contents).Decode(&updatedUsr)
    // Append coordinate to the user object
    err = getCoordinates(&updatedUsr)
    if err != nil {
        return models.User{}, err
    }

    // Grab id
    objId := bson.ObjectIdHex(id)
    // make connection
    conn := uc.session.DB("cmpe273_asgmt2").C("user")
    err = conn.Update(bson.M{"id": objId}, updatedUsr)
    if err != nil {
        log.Println(err)
        return models.User{}, err
    }

    return updatedUsr, nil
}
