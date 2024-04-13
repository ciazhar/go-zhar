package service

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/crud-testcontainers/internal/person/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/crud-testcontainers/internal/person/repository"
)

type PersonService interface {
	Insert(person *model.Person) error
	InsertBatch(persons *[]model.Person) error
	FindOne(id string) (model.Person, error)
	FindAll(
		name string,
		email string,
		age int,
	) ([]model.Person, error)

	FindCountry(country string) ([]model.Person, error)
	FindAgeRange(min int, max int) ([]model.Person, error)
	FindHobby(hobby []string) ([]model.Person, error)
	FindMinified() ([]model.PersonMinified, error)
	Update(id string, person model.UpdatePersonForm) error
	Delete(id string) error
}

type personService struct {
	p repository.PersonRepository
}

func (p personService) Insert(person *model.Person) error {
	return p.p.Insert(person)
}

func (p personService) InsertBatch(persons *[]model.Person) error {
	return p.p.InsertBatch(persons)
}

func (p personService) FindOne(id string) (model.Person, error) {
	return p.p.FindOne(id)
}

func (p personService) FindAll(name string, email string, age int) ([]model.Person, error) {

	return p.p.FindAll(name, email, age)
}

func (p personService) FindCountry(country string) ([]model.Person, error) {

	return p.p.FindCountry(country)
}

func (p personService) FindAgeRange(min int, max int) ([]model.Person, error) {

	return p.p.FindAgeRange(min, max)
}

func (p personService) FindHobby(hobby []string) ([]model.Person, error) {

	return p.p.FindHobby(hobby)
}

func (p personService) FindMinified() ([]model.PersonMinified, error) {

	return p.p.FindMinified()
}

func (p personService) Update(id string, person model.UpdatePersonForm) error {

	return p.p.Update(id, person)
}

func (p personService) Delete(id string) error {

	return p.p.Delete(id)
}

func NewPersonService(p repository.PersonRepository) PersonService {
	return personService{
		p: p,
	}
}
