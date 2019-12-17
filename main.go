package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// LogRequest collection struct to save data to the database
type LogRequest struct {
	User      string `json:"user_id"`
	Operation string `json:"operation"`
	Entity    string `json:"entity,omitempty"`
	Message   string `json:"message,omitempty"`
}

func index(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Write([]byte(`{"message": "Welcome to activity logger microservices"}`))
}

// LogCreate controller for creating new Logs
func LogCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	connString := os.Getenv("DB")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		InternalError(w)
		return
	}
	defer db.Close()
	var log LogRequest
	json.NewDecoder(r.Body).Decode(&log)
	createdat := time.Now().Format("2006-01-02 15:04:05")
	var id int
	err = db.QueryRow(`INSERT INTO logs (user_id, operation, entity, message, created_at) 
	VALUES('` + log.User + `','` + log.Operation + `', '` + log.Entity + `', '` + log.Message + `', to_timestamp('` + createdat + `', 'YYYY-MM-DD HH24:MI:SS')) RETURNING id`).Scan(&id)
	if err != nil {
		InternalError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "Logged successfully." }`))
}

// func LogByUser(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Set("Content-Type", "application/json")
// 	var logg []*Log
// 	params := mux.Vars(req)
// 	id := params["id"]
// 	// page := params["page"]
// 	// pg, _ := strconv.ParseInt(page,10,64)
// 	number := params["no"]
// 	no, _ := strconv.ParseInt(number, 10, 64)
// 	findOptions := options.Find()
// 	findOptions.SetLimit(no)
// 	// skip := no*pg
// 	ObjId, _ := primitive.ObjectIDFromHex(id)
// 	query := bson.M{"user_id": ObjId}
// 	collection := database.Collection("logs")
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
// 	cur, err := collection.Find(ctx, query, findOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for cur.Next(ctx) {
// 		var elem Log
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		logg = append(logg, &elem)
// 	}
// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	cur.Close(ctx)
// 	json.NewEncoder(res).Encode(logg)
// }

func handleRequests() {
	fmt.Printf("Activity logger started 🔥🔥🔥")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/", LogCreate).Methods("POST")
	// router.HandleFunc("/{id}/{page}/{no}", LogByUser).Methods("GET")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}

func main() {
	handleRequests()
}

// InternalError - emits internal server error
func InternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "Something went wrong."}`))
}
