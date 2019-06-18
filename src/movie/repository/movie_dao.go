package repository

import (
	"github.com/ciazhar/go-mongo-example/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type MovieDao struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "movies"
)

// Establish a connection to database
func (m *MovieDao) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of movies
func (m *MovieDao) FindAllMovieByQueryAndPaged(q interface{}, skip int, limit int) ([]model.Movie, error) {
	var movie []model.Movie
	err := db.C(COLLECTION).Find(q).Skip((skip - 1) * 10).Limit(limit).All(&movie)
	return movie, err
}

// Find a movie by its id
func (m *MovieDao) FindById(id string) (model.Movie, error) {
	var movie model.Movie
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&movie)
	return movie, err
}

// Insert a movie into database
func (m *MovieDao) Insert(movie model.Movie) error {
	err := db.C(COLLECTION).Insert(&movie)
	return err
}

// Delete an existing movie
func (m *MovieDao) Delete(movie model.Movie) error {
	err := db.C(COLLECTION).Remove(&movie)
	return err
}

// Update an existing movie
func (m *MovieDao) Update(movie model.Movie) error {
	err := db.C(COLLECTION).UpdateId(movie.ID, &movie)
	return err
}
