package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/model"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/repository"
)

type BookService struct {
	bookRepository *repository.BookRepository
}

func (b *BookService) Insert(context context.Context, book *model.Book) error {
	return b.bookRepository.Insert(context, book)
}

func NewBookService(bookRepository *repository.BookRepository) *BookService {

	return &BookService{
		bookRepository: bookRepository,
	}
}
