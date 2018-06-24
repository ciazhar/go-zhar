package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"ciazhar.com/go-mongo-example/controller"
	"ciazhar.com/go-mongo-example/config"
	"ciazhar.com/go-mongo-example/dao"
)

var conf = config.Config{}
var movieDao = dao.MovieDao{}

func init() {
	conf.Read()

	movieDao.Server = conf.Server
	movieDao.Database = conf.Database
	movieDao.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", controller.AllMoviesByQueryAndPagedEndpoint).Methods("GET")
	r.HandleFunc("/movies", controller.CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", controller.UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", controller.DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", controller.FindMovieEndpoint).Methods("GET")

	handler:= cors.Default().Handler(r)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}