package main

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrBookAlreadyExists = errors.New("Book with this ISBN already exists")
	ErrBookNotFound      = errors.New("Book not found")
)

func getAllBooks() ([]Book, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(DB).C(BOOKS_COLLECTION)
	var books []Book
	err := c.Find(bson.M{}).All(&books)
	return books, err
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
	s := session.Copy()
	defer s.Close()

	c := s.DB(DB).C(BOOKS_COLLECTION)
	book := new(Book)
	err := c.Find(bson.M{"isbn": isbn}).One(book)
	if err != nil {
		return nil, err
	}
	if book.ISBN == "" {
		return nil, ErrBookNotFound
	}
	return book, nil
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
