package routes

import (
	"Backend/controllers"
	"github.com/gorilla/mux"
)

func ComandoRoute(router *mux.Router) {
	//All routes related to users comes here
	router.HandleFunc("/", Controllers.IndexHandler).Methods("GET") //add this
	router.HandleFunc("/Principal", Controllers.RequestPrincipal()) //.Methods("GET") //add this
	router.HandleFunc("/Kill", Controllers.RequestKill())           //.Methods("GET")
	router.HandleFunc("/Cpu", Controllers.RequestCPU())             //.Methods("GET")
	router.HandleFunc("/Memoria", Controllers.RequestMemory())      //.Methods("GET")
}
