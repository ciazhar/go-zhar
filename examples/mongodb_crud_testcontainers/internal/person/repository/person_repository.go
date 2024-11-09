package repository

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
	mongo2 "github.com/ciazhar/go-start-small/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonRepository struct {
	c *mongo.Collection
}

func (p *PersonRepository) Delete(ctx context.Context, id string) (err error) {

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	_, err = p.c.DeleteOne(ctx, bson.D{{"_id", hex}})
	return
}

func (p *PersonRepository) Update(ctx context.Context, id string, person model.UpdatePersonForm) (err error) {

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	_, err = p.c.UpdateOne(ctx, bson.D{{"_id", hex}}, bson.D{{"$set", person}})
	return
}

func (p *PersonRepository) FindMinified(ctx context.Context) (res []model.PersonMinified, err error) {

	projection := bson.M{
		"_id":     1,
		"name":    1,
		"country": "$address.country",
	}

	cursor, err := p.c.Find(ctx, bson.D{}, options.Find().SetProjection(projection))
	if err != nil {
		return
	}
	if err = cursor.All(ctx, &res); err != nil {
		return
	}

	return
}

func (p *PersonRepository) FindHobby(ctx context.Context, hobby []string) (res []model.Person, err error) {

	filter := bson.D{}
	if len(hobby) != 0 {
		filter = append(filter, bson.E{Key: "hobbies", Value: bson.D{{"$in", hobby}}})
	}
	cursor, err := p.c.Find(ctx, filter)
	if err != nil {
		return
	}
	if err = cursor.All(ctx, &res); err != nil {
		return
	}
	return
}

func (p *PersonRepository) FindAgeRange(ctx context.Context, min int, max int) (res []model.Person, err error) {

	filter := bson.D{}
	if min != 0 {
		filter = append(filter, bson.E{Key: "age", Value: bson.D{{"$gte", min}}})
	}
	if max != 0 {
		filter = append(filter, bson.E{Key: "age", Value: bson.D{{"$lte", max}}})
	}
	cursor, err := p.c.Find(ctx, filter)
	if err != nil {
		return
	}
	if err = cursor.All(ctx, &res); err != nil {
		return
	}
	return
}

func (p *PersonRepository) FindCountry(ctx context.Context, country string) (res []model.Person, err error) {

	filter := bson.D{}
	if country != "" {
		filter = append(filter, bson.E{Key: "address.country", Value: country})
	}
	cursor, err := p.c.Find(ctx, filter)
	if err != nil {
		return
	}
	if err = cursor.All(ctx, &res); err != nil {
		return
	}
	return
}

func (p *PersonRepository) FindAllPageSize(
	ctx context.Context,
	page, size int,
	sort string, // format "param1,asc;param2,desc"
	name string,
	email string,
	age int,
) (res []model.Person, err error) {

	findOptions, err := mongo2.CreatePagingAndSortingOptions(page, size, sort)
	if err != nil {
		return
	}

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

	logger.LogInfo(ctx, "filter", map[string]interface{}{"filter": filter})

	c, err := p.c.Find(ctx, filter, findOptions)
	if err != nil {
		return
	}
	if err = c.All(ctx, &res); err != nil {
		return
	}
	return
}

func (p *PersonRepository) FindOne(ctx context.Context, id string) (res model.Person, err error) {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = p.c.FindOne(ctx, bson.M{"_id": hex}).Decode(&res)
	return
}

func (p *PersonRepository) Insert(ctx context.Context, person *model.Person) (err error) {
	res, err := p.c.InsertOne(ctx, person)
	if err != nil {
		return
	}
	person.Id = res.InsertedID.(primitive.ObjectID)
	return
}

func (p *PersonRepository) InsertBatch(ctx context.Context, persons *[]model.Person) (err error) {
	var interfaceUsers []interface{}

	for i := range *persons {
		interfaceUsers = append(interfaceUsers, (*persons)[i])
	}
	res, err := p.c.InsertMany(ctx, interfaceUsers)
	if res != nil {
		for i := range *persons {
			(*persons)[i].Id = res.InsertedIDs[i].(primitive.ObjectID)
		}
	}
	return
}

func NewPersonRepository(c *mongo.Database) *PersonRepository {
	return &PersonRepository{c: c.Collection("persons")}
}
