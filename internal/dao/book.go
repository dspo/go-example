package dao

import (
	"log"

	"gorm.io/gorm"

	"gitee.com/huajinet/go-example/internal/model"
)

type BookDao interface {
	CreateBook(book *model.Book) error
	List() (int64, []*model.Book, error)
}

func NewBook(db *gorm.DB) BookDao {
	return &bookDao{db: db.Debug()}
}

type bookDao struct {
	db *gorm.DB
}

func (b bookDao) CreateBook(book *model.Book) error {
	return b.db.Model(book).Create(book).Error
}

func (b bookDao) List() (total int64, list []*model.Book, err error) {
	err = b.db.Model(new(model.Book)).
		Count(&total).
		Find(&list).Error
	if err != nil {
		log.Printf("failed to find books: %v\n", err)
	}
	return
}
