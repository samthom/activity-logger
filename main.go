package main

import (
	"log"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"github.com/gorilla/mux"
	"time"
)

// var client *mongo.Client
var database *mongo.Database

// log collection struct to save data to the database
type Log struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User primitive.ObjectID `json:"user_id" bson:"user_id"`
	Operation string `json:"operation,omitempty" bson:"operation,omitempty"`
	Entity string `json:"entity,omitempty" bson:"entity,omitempty"`
	Info string `json:"info,omitempty" bson:"info,omitempty"`
	Created string `json:"created_at" bson:"created_at"`
}

/*
 *--------------------------------------------------------------
 * Index Controller (Homepage)
 *--------------------------------------------------------------
 * Controller for showing the homepage of the service
 * @TODO refactoring pending
 */
func index(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Write([]byte(`{"message": "Welcome to activity logger microservices"}`))
}

/*
 *---------------------------------------------------------------
 * Create Create operation Log Controller
 *---------------------------------------------------------------
 * @TODO refactoring pending
 * Controller for creating a login log
 *
 */
func LogCreate(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var log Log
	json.NewDecoder(req.Body).Decode(&log)
	collection := database.Collection("logs")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, log)
	json.NewEncoder(res).Encode(result)
}
/*
 *------------------------------------------------------------
 * Get logs by user
 *------------------------------------------------------------
 * Controller for getting the logs of the user
 * according to the user name sent by the microservice
 * paginating the logs of the user
 */
func LogByUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var logg []*Log
	params := mux.Vars(req)
	id := params["id"]
	// page := params["page"]
	number := params["no"]
	no, _ := strconv.ParseInt(number,10,64)
	findOptions := options.Find()
	findOptions.SetLimit(no)

	ObjId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"user_id": ObjId}
	collection := database.Collection("logs")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, query, findOptions)
	if (err != nil) {
		log.Fatal(err)
	}
	for  cur.Next(ctx) {
		var elem Log
		err := cur.Decode(&elem)
		if (err != nil) {
			log.Fatal(err)
		}
		logg = append(logg, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(ctx)
	json.NewEncoder(res).Encode(logg)
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
	fmt.Printf("Activity logger started 🔥🔥🔥")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/", LogCreate).Methods("POST")
	router.HandleFunc("/{id}/{page}/{no}", LogByUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
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
	// Client contest is created in the main func and it will be accessible for the entire program
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	database = client.Database("kjc")
	if (err != nil) {
		log.Fatal(err)
	}
	handleRequests()
}

