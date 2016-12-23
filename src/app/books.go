package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"app/cache"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrBookAlreadyExists = errors.New("Book with this ISBN already exists")
	ErrBookNotFound      = errors.New("Book not found")
)

const (
	ALL_BOOKS    = "all_books"
	BOOK_BY_ISBN = "book_by_isbn"
)

func getAllBooks() ([]Book, error) {
	cachedBooks, _ := cache.Get(ALL_BOOKS)
	var books []Book
	err := json.NewDecoder(strings.NewReader(cachedBooks)).Decode(&books)
	if err != nil {
		log.Println("Using the database")
		s := session.Copy()
		defer s.Close()

		c := s.DB(DB).C(BOOKS_COLLECTION)
		err = c.Find(bson.M{}).All(&books)
		updateCache(ALL_BOOKS, books...)
	}
	return books, err
}

func updateCache(key string, books ...Book) {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(books)
	cache.Set(key, buffer.String())
}

func saveBook(book *Book) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(DB).C(BOOKS_COLLECTION)
	err := c.Insert(book)
	if err != nil {
		if mgo.IsDup(err) {
			return ErrBookAlreadyExists
		}
	}
	return err
}

func getBookByIsbn(isbn string) (*Book, error) {
	cachedBook, _ := cache.Get(BOOK_BY_ISBN)
	book := new(Book)
	err := json.NewDecoder(strings.NewReader(cachedBook)).Decode(book)
	if err != nil {
		log.Println("Using the database")
		s := session.Copy()
		defer s.Close()

		c := s.DB(DB).C(BOOKS_COLLECTION)
		err := c.Find(bson.M{"isbn": isbn}).One(book)
		if err != nil {
			return nil, err
		}
		if book.ISBN == "" {
			return nil, ErrBookNotFound
		}
		updateCache(BOOK_BY_ISBN, *book)
		return book, nil
	}
	return book, err
}

func updateBook(isbn string, book *Book) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(DB).C(BOOKS_COLLECTION)
	err := c.Update(bson.M{"isbn": isbn}, book)
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrBookNotFound
		}
		return err
	}
	return nil
}

func removeBook(isbn string) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(DB).C(BOOKS_COLLECTION)
	err := c.Remove(bson.M{"isbn": isbn})
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrBookNotFound
		}
		return err
	}
	return nil
}

func ClearDatabase() {
	s := session.Copy()
	defer s.Close()
	c := s.DB(DB).C(BOOKS_COLLECTION)
	c.RemoveAll(bson.M{})
}
