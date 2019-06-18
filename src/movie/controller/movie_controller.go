package controller

import (
	"encoding/json"
	"github.com/ciazhar/go-mongo-example/model"
	"github.com/ciazhar/go-mongo-example/movie/repository"
	"github.com/ciazhar/go-mongo-example/tool"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

var movieDao = dao.MovieDao{}

// GET list of movies
func AllMoviesByQueryAndPagedEndpoint(w http.ResponseWriter, r *http.Request) {
	q := map[string]interface{}{}

	name := r.URL.Query().Get("name")
	coverImage := r.URL.Query().Get("coverImage")
	description := r.URL.Query().Get("description")
	skip, _ := strconv.Atoi(r.URL.Query().Get("skip"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if name != "" {
		q["name"] = name
	}
	if coverImage != "" {
		q["cover_image"] = coverImage
	}
	if description != "" {
		q["description"] = description
	}
	if skip == 0 {
		skip = 1
	}
	if limit == 0 {
		limit = 20
	}

	movies, err := movieDao.FindAllMovieByQueryAndPaged(q, skip, limit)
	if err != nil {
		tool.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tool.RespondWithJson(w, http.StatusOK, movies)
}

// GET a movie by its ID
func FindMovieEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := movieDao.FindById(params["id"])
	if err != nil {
		tool.RespondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}
	tool.RespondWithJson(w, http.StatusOK, movie)
}

// POST a new movie
func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		tool.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.ID = bson.NewObjectId()
	if err := movieDao.Insert(movie); err != nil {
		tool.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tool.RespondWithJson(w, http.StatusCreated, movie)
}

// PUT update an existing movie
func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		tool.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := movieDao.Update(movie); err != nil {
		tool.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tool.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing movie
func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		tool.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := movieDao.Delete(movie); err != nil {
		tool.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tool.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
