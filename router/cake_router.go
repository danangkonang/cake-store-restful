package router

import (
	"github.com/danangkonang/cake-store-restful/config"
	"github.com/danangkonang/cake-store-restful/controller"
	"github.com/danangkonang/cake-store-restful/service"
	"github.com/gorilla/mux"
)

func CakeRouter(router *mux.Router, db *config.DB) {
	c := controller.NewCakeController(
		service.NewServiceCake(db),
	)
	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/cakes", c.FindCakes).Methods("GET")
	v1.HandleFunc("/cake", c.SaveCake).Methods("POST")
	v1.HandleFunc("/cake/{id:[0-9]+}", c.FindCakeById).Methods("GET")
	v1.HandleFunc("/cake", c.UpdateCake).Methods("PUT")
	v1.HandleFunc("/cake", c.DeleteCake).Methods("DELETE")
}
