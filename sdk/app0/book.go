package app0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitee.com/huajinet/go-example/internal/model"
)

type Book interface {
	Create(book *model.Book) error
	List() (uint64, []*model.Book, error)
}

func NewBook() Book {
	return &book{}
}

type book struct {
}

func (b book) Create(book *model.Book) error {
	var buf = bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(book); err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, "http://app0:80/api/v1/books", buf)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("create book failed with status code: %v", response.StatusCode)
	}

	return json.NewDecoder(response.Body).Decode(book)
}

func (b book) List() (uint64, []*model.Book, error) {
	response, err := http.Get("http://app0:80/api/v1/books")
	if err != nil {
		return 0, nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return 0, nil, fmt.Errorf("list book failed with status code: %v", response.StatusCode)
	}

	type Data struct {
		Total uint64        `json:"total"`
		List  []*model.Book `json:"list"`
	}

	var data Data
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0, nil, err
	}

	return data.Total, data.List, nil
}
