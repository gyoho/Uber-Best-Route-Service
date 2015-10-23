package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    "fmt"
    "log"
    "gopkg.in/mgo.v2"
    "./controllers"
)

func main() {
    // Instantiate a new router
    router := httprouter.New()

    // Get a UserController instance
    userController := controllers.NewUserController(getMongoSession())

    // Get a user resource
    router.GET("/locations/:id", userController.GetUser)
    router.POST("/locations", userController.CreateUser)
    router.PUT("/locations/:id", userController.UpdateUser)
    router.DELETE("/locations/:id", userController.RemoveUser)

    // Fire up the server
    fmt.Println("Server listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getMongoSession() *mgo.Session {
    // Connect to our local mongo
    session, err := mgo.Dial("mongodb://localhost")
    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }

    return session
}
