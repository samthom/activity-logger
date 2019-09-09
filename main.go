package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	// "labix.org/v2/mgo/bson"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"github.com/gorilla/mux"
	"time"
)

// log collection struct to save data to the database
type Log struct {
	ID primitive.ObjectID `json: "_id,omitempty" bson: "_id,omitempty"`
	User primitive.ObjectID `json: "user_id,omitempty" bson: "user_id,omitempty"`
	Operation string `json:"operation,omitempty" bson: "operation,omitempty"`
	Entity string `json:"entity,omitempty" bson: "entity,omitempty"`
	Created string `json:"created_at" bson: "created_at"`
}

/*
 *---------------------------------------------------------------------
 * Main Func
 *---------------------------------------------------------------------
 * Main function starts from here. Have to call the routes to this
 * function for execution
 *
 */
func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if (err != nil) {
		log.Fatal(err)
	}
	collection := client.Database("testing").Collection("test")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{ "name": "pi", "value": 3.141 })
	id := res.InsertedID
	fmt.Printf("Updated Database %s", id)
	handleRequests()
}

/*
 *----------------------------------------------------------------
 * Routes func
 *----------------------------------------------------------------
 * @TODO Have refactor the code into better file structure
 * @TODO Have to implement authentication middleware
 * routes are composed according to the operations in creating logs
 *
 */
func handleRequests() {
	fmt.Printf("Activity logger started")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/login", PostLogin).Methods("POST")
	router.HandleFunc("/create", PostCreate).Methods("POST")
	router.HandleFunc("/update", PostUpdate).Methods("POST")
	router.HandleFunc("/delete", PostDelete).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

/*
 *--------------------------------------------------------------
 * Index Controller (Homepage)
 *--------------------------------------------------------------
 * Controller for showing the homepage of the service
 * @TODO refactoring pending
 */
func index(res http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res.Write([]byte(`{"message": "Welcome to activity logger microservices"}`))
}

/*
 *---------------------------------------------------------------
 * Create Login Log Controller
 *---------------------------------------------------------------
 * @TODO refactoring pending
 * Controller for creating a login log
 *
 */
func PostLogin(w http.ResponseWriter, req *http.Request) {
	// mux.Vars for getting the values from the url not needed here
	// params := mux.Vars(req)
	var log Log
	// Decoding the request into the struct
	// But for inserting into the database have to use bson
	_ = json.NewDecoder(req.body).Decode(&log)
}

/*
 *---------------------------------------------------------------
 * Create Create operation Log Controller
 *---------------------------------------------------------------
 * @TODO refactoring pending
 * Controller for creating a login log
 *
 */
func PostCreate(w http.ResponseWriter, req *http.Request) {

}

/*
 *---------------------------------------------------------------
 * Create Update operation Log Controller
 *---------------------------------------------------------------
 * @TODO refactoring pending
 * Controller for creating a login log
 *
 */
func PostUpdate(w http.ResponseWriter, req *http.Request) {

}

/*
 *---------------------------------------------------------------
 * Create Delete operation Log Controller
 *---------------------------------------------------------------
 * @TODO refactoring pending
 * Controller for creating a login log
 *
 */
func PostDelete(w http.ResponseWriter, req *http.Request) {

}
