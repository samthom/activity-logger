/**
TODO
1. Have to refactor the db connecting call into another function
2. Have to refactor the response creating part of the controller
**/

package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/samthom/activity-logger/controller/errors"
)

// LogRequest collection struct to save data to the database
type LogRequest struct {
	ID        int       `json:"id"`
	User      string    `json:"user_id"`
	Operation string    `json:"operation"`
	Entity    string    `json:"entity,omitempty"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// LogResponse collection struct to save data to the database
type LogResponse struct {
	Data []LogRequest `json:"data"`
}

// Index controller for pinging
func Index(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Write([]byte(`{"message": "Welcome to activity logger microservices"}`))
}

// LogCreate controller for creating new Logs
func LogCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	connString := os.Getenv("DB")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		errors.InternalError(w)
		return
	}
	defer db.Close()
	var log LogRequest
	json.NewDecoder(r.Body).Decode(&log)
	createdat := time.Now().Format("2006-01-02 15:04:05")
	var id int
	err = db.QueryRow(`INSERT INTO logs (user_id, operation, entity, message, created_at) 
	VALUES('` + log.User + `','` + log.Operation + `', '` + log.Entity + `', '` + log.Message + `', to_timestamp('` + createdat + `', 'YYYY-MM-DD HH24:MI:SS')) 
	RETURNING id`).Scan(&id)
	if err != nil {
		errors.InternalError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "Logged successfully." }`))
}

// LogByUser return the logs in a paginated way
func LogByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	connString := os.Getenv("DB")
	db, err := sql.Open("postgres", connString)
	var response LogResponse
	if err != nil {
		errors.InternalError(w)
		return
	}
	defer db.Close()
	params := mux.Vars(r)
	id := params["id"]
	page := params["page"]
	p, err := strconv.Atoi(page)
	no := params["no"]
	n, err := strconv.Atoi(no)
	offset := n * p
	if p == 1 {
		offset = 0
	}
	o := strconv.Itoa(offset)
	qry := `SELECT * FROM logs WHERE user_id=$1 LIMIT $2 OFFSET $3;`
	iterator, err := db.Query(qry, id, no, o)
	if err != nil {
		errors.InternalError(w)
		return
	}
	response, err = logIterator(iterator)
	if err != nil {
		errors.InternalError(w)
		return
	}
	json.NewEncoder(w).Encode(response)
}

// GetLogsPaginated - Controller for returning the logs paginated
func GetLogsPaginated(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response LogResponse
	connString := os.Getenv("DB")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		errors.InternalError(w)
		return
	}
	defer db.Close()
	params := mux.Vars(r)
	page := params["page"]
	no := params["no"]
	p, err := strconv.Atoi(page)
	n, err := strconv.Atoi(no)
	offset := n * p
	if p == 1 {
		offset = 0
	}
	o := strconv.Itoa(offset)
	qry := `SELECT * FROM logs LIMIT $1 OFFSET $2;`
	iterator, err := db.Query(qry, no, o)
	if err != nil {
		errors.InternalError(w)
		return
	}
	response, err = logIterator(iterator)
	if err != nil {
		errors.InternalError(w)
		return
	}
	json.NewEncoder(w).Encode(response)
}

// LogsByEntity - Controller for returning the logs according to the provided entity
func LogsByEntity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response LogResponse
	connStr := os.Getenv("DB")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		errors.InternalError(w)
		return
	}
	defer db.Close()
	params := mux.Vars(r)
	entity := params["entity"]
	page := params["page"]
	no := params["no"]
	p, err := strconv.Atoi(page)
	n, err := strconv.Atoi(no)
	offset := p * n
	if p == 1 {
		offset = 0
	}
	o := strconv.Itoa(offset)
	qry := `SELECT * FROM logs WHERE entity=$1 LIMIT $2 OFFSET $3`
	iterator, err := db.Query(qry, entity, no, o)
	response, err = logIterator(iterator)
	if err != nil {
		errors.InternalError(w)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func logIterator(iterator *sql.Rows) (LogResponse, error) {
	defer iterator.Close()
	var response LogResponse
	for iterator.Next() {
		var log LogRequest
		err := iterator.Scan(
			&log.ID,
			&log.User,
			&log.Operation,
			&log.Entity,
			&log.Message,
			&log.CreatedAt,
		)
		if err != nil {
			return response, err
		}
		response.Data = append(response.Data, log)
	}
	err := iterator.Err()
	if err != nil {
		return response, err
	}
	return response, nil
}
