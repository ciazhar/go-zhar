package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/crud-testcontainers/internal/person/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonRepository interface {
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

type personRepository struct {
	c *mongo.Collection
}

func (p personRepository) Delete(id string) error {

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = p.c.DeleteOne(context.Background(), bson.D{{"_id", hex}})
	return err
}

func (p personRepository) Update(id string, person model.UpdatePersonForm) error {

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = p.c.UpdateOne(context.Background(), bson.D{{"_id", hex}}, bson.D{{"$set", person}})
	return err
}

func (p personRepository) FindMinified() ([]model.PersonMinified, error) {

	projection := bson.M{
		"_id":     1,
		"name":    1,
		"country": "$address.country",
	}

	cursor, err := p.c.Find(context.Background(), bson.D{}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	var results []model.PersonMinified
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (p personRepository) FindHobby(hobby []string) ([]model.Person, error) {

	filter := bson.D{}
	if len(hobby) != 0 {
		filter = append(filter, bson.E{Key: "hobbies", Value: bson.D{{"$in", hobby}}})
	}
	cursor, err := p.c.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var results []model.Person
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (p personRepository) FindAgeRange(min int, max int) ([]model.Person, error) {

	filter := bson.D{}
	if min != 0 {
		filter = append(filter, bson.E{Key: "age", Value: bson.D{{"$gte", min}}})
	}
	if max != 0 {
		filter = append(filter, bson.E{Key: "age", Value: bson.D{{"$lte", max}}})
	}
	cursor, err := p.c.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var results []model.Person
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (p personRepository) FindCountry(country string) ([]model.Person, error) {

	filter := bson.D{}
	if country != "" {
		filter = append(filter, bson.E{Key: "address.country", Value: country})
	}
	cursor, err := p.c.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var results []model.Person
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (p personRepository) FindAll(
	name string,
	email string,
	age int,
) ([]model.Person, error) {

	filter := bson.D{}
	if name != "" {
		filter = append(filter, bson.E{Key: "name", Value: name})
	}
	if email != "" {
		filter = append(filter, bson.E{Key: "email", Value: email})
	}
	if age != 0 {
		filter = append(filter, bson.E{Key: "age", Value: age})
	}
	cursor, err := p.c.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var results []model.Person
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (p personRepository) FindOne(id string) (person model.Person, err error) {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Person{}, err
	}
	err = p.c.FindOne(context.Background(), bson.M{"_id": hex}).Decode(&person)
	return person, err
}

func (p personRepository) Insert(person *model.Person) error {
	res, err := p.c.InsertOne(context.Background(), person)
	if err != nil {
		return err
	}
	person.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (p personRepository) InsertBatch(persons *[]model.Person) error {
	var interfaceUsers []interface{}

	for i := range *persons {
		interfaceUsers = append(interfaceUsers, (*persons)[i])
	}
	res, err := p.c.InsertMany(context.Background(), interfaceUsers)
	if res != nil {
		for i := range *persons {
			(*persons)[i].Id = res.InsertedIDs[i].(primitive.ObjectID)
		}
	}
	return err
}

func NewPersonRepository(c *mongo.Database) PersonRepository {
	return &personRepository{c: c.Collection("persons")}

}
