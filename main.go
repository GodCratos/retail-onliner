package main

import (
	"net/http"
	"runtime"

	"github.com/GodCratos/retail-onliner/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/update_order", handlers.HandlerUpdateOrder).Methods("POST")
	http.Handle("/", router)
	go http.ListenAndServe(":8181", nil)
	runtime.Goexit()
}
