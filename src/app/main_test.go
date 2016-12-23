package main_test

import (
	main "app"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	a := assert.New(t)
	setup()

	body := `{
    "isbn":    "0134190440",
    "title":   "The Go Programming Language",
    "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
    "price":   "$34.57"
  }`

	w := createBook(body)

	a.Equal(http.StatusCreated, w.Code)
}

func TestGetBookByIsbn(t *testing.T) {
	a := assert.New(t)
	setup()
	body := `{
    "isbn":    "0134190440",
    "title":   "The Go Programming Language",
    "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
    "price":   "$34.57"
  }`
	createBook(body)

	r := main.NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/0134190440", nil)
	r.ServeHTTP(w, req)

	log.Println(w.Body.String())
	a.Equal(http.StatusOK, w.Code)
	book := new(main.Book)
	json.NewDecoder(w.Body).Decode(book)
	a.Equal("The Go Programming Language", book.Title)
	a.Equal(2, len(book.Authors))
}

func TestDeleteBook(t *testing.T) {
	a := assert.New(t)
	setup()
	body := `{
    "isbn":    "0134190440",
    "title":   "The Go Programming Language",
    "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
    "price":   "$34.57"
  }`
	createBook(body)

	r := main.NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/0134190440", nil)
	r.ServeHTTP(w, req)

	log.Println(w.Body.String())
	a.Equal(http.StatusNoContent, w.Code)
}

func TestAllBooks(t *testing.T) {
	a := assert.New(t)
	setup()
	body := `{
    "isbn":    "0134190440",
    "title":   "The Go Programming Language",
    "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
    "price":   "$34.57"
  }`
	createBook(body)
	body = `{
    "isbn":    "0134190441",
    "title":   "Advanced Go Programming",
    "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
    "price":   "$40.57"
  }`
	createBook(body)

	r := main.NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)

	log.Println(w.Body.String())
	a.Equal(http.StatusOK, w.Code)
	books := []main.Book{}
	json.NewDecoder(w.Body).Decode(&books)
	a.Equal(2, len(books))
}

func createBook(body string) *httptest.ResponseRecorder {
	router := main.NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

func setup() {
	main.ClearDatabase()
}
