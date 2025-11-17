package service

import (
	"gitee.com/huajinet/go-example/internal/dao"
	"gitee.com/huajinet/go-example/internal/model"
)

type Book interface {
	Create(book *model.Book) error
	List() (int64, []*model.Book, error)
}

func NewBook(da dao.BookDao) Book {
	return &book{da: da}
}

type book struct {
	da dao.BookDao
}

func (b book) Create(book *model.Book) error {
	return b.da.CreateBook(book)
}

func (b book) List() (int64, []*model.Book, error) {
	return b.da.List()
}
