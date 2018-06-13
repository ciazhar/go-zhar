package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	."go-mongo-example/dao"
	."go-mongo-example/model"
	."go-mongo-example/tool"
	"strconv"
)

var dao = MovieDao{}

// GET list of movies
func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
	var limit =20
	var skip  =1

	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	movies, err := dao.FindAll(limit,skip)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusOK, movies)
}

// GET a movie by its ID
func FindMovieEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := dao.FindById(params["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}
	RespondWithJson(w, http.StatusOK, movie)
}

// POST a new movie
func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.ID = bson.NewObjectId()
	if err := dao.Insert(movie); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusCreated, movie)
}

// PUT update an existing movie
func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(movie); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing movie
func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(movie); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}