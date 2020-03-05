package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/samthom/activity-logger/controller"
)

// HandleRequests - Routing function for the application and starts the serever
func HandleRequests() {
	fmt.Printf("Activity logger started 🔥🔥🔥")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", controller.Index).Methods("GET")
	router.HandleFunc("/", controller.LogCreate).Methods("POST")
	router.HandleFunc("/user/{id}/{page}/{no}", controller.LogByUser).Methods("GET")
	router.HandleFunc("/{page}/{no}", controller.GetLogsPaginated).Methods("GET")
	router.HandleFunc("/entity/{entity}/{page}/{no}", controller.LogsByEntity).Methods("GET")
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
