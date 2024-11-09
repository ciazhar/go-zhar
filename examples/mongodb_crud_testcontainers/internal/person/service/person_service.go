package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/model"
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/repository"
)

type PersonService struct {
	p *repository.PersonRepository
}

func (p *PersonService) Insert(ctx context.Context, person *model.Person) error {
	return p.p.Insert(ctx, person)
}

func (p *PersonService) InsertBatch(ctx context.Context, persons *[]model.Person) error {
	return p.p.InsertBatch(ctx, persons)
}

func (p *PersonService) FindOne(ctx context.Context, id string) (model.Person, error) {
	return p.p.FindOne(ctx, id)
}

func (p *PersonService) FindAllPageSize(ctx context.Context,
	page, size int,
	sort string,
	name string,
	email string,
	age int,
) ([]model.Person, error) {
	return p.p.FindAllPageSize(ctx, page, size, sort, name, email, age)
}

func (p *PersonService) FindCountry(ctx context.Context, country string) ([]model.Person, error) {
	return p.p.FindCountry(ctx, country)
}

func (p *PersonService) FindAgeRange(ctx context.Context, min int, max int) ([]model.Person, error) {
	return p.p.FindAgeRange(ctx, min, max)
}

func (p *PersonService) FindHobby(ctx context.Context, hobby []string) ([]model.Person, error) {
	return p.p.FindHobby(ctx, hobby)
}

func (p *PersonService) FindMinified(ctx context.Context) ([]model.PersonMinified, error) {
	return p.p.FindMinified(ctx)
}

func (p *PersonService) Update(ctx context.Context, id string, person model.UpdatePersonForm) error {
	return p.p.Update(ctx, id, person)
}

func (p *PersonService) Delete(ctx context.Context, id string) error {
	return p.p.Delete(ctx, id)
}

func NewPersonService(p *repository.PersonRepository) *PersonService {
	return &PersonService{
		p: p,
	}
}
