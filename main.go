package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	. "go-mongo-example/controller"
	. "go-mongo-example/dao"
	. "go-mongo-example/config"
	"github.com/rs/cors"
)

var config = Config{}
var movieDao = MovieDao{}

func init() {
	config.Read()

	movieDao.Server = config.Server
	movieDao.Database = config.Database
	movieDao.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", FindMovieEndpoint).Methods("GET")

	handler:= cors.Default().Handler(r)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}