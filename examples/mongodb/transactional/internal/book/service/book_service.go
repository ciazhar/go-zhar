package service

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/book/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/book/repository"
)

type BookService interface {
	Insert(context context.Context, book *model.Book) error
}

type bookService struct {
	bookRepository repository.BookRepository
}

func (b bookService) Insert(context context.Context, book *model.Book) error {
	return b.bookRepository.Insert(context, book)
}

func NewBookService(bookRepository repository.BookRepository) BookService {

	return &bookService{
		bookRepository: bookRepository,
	}
}
