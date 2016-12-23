package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
)

const (
	DB               = "bookStore"
	BOOKS_COLLECTION = "books"
)

type Book struct {
	ISBN    string   `json:"isbn"`
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Price   string   `json:"price"`
}

var session *mgo.Session
var err error

func init() {
	server := os.Getenv("MONGO_SERVER_URL")
	if server == "" {
		server = "localhost"
	}
	log.Println("Connecting to mongo on: " + server + ":27017")
	session, err = mgo.Dial(server + ":27017")
	if err != nil {
		panic(err)
	}
}

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/books", logRequest(allBooks))
	router.POST("/books", logRequest(addBook))
	router.GET("/books/:isbn", logRequest(booksByIsbn))
	router.PUT("/books/:isbn", logRequest(editBook))
	router.DELETE("/books/:isbn", logRequest(deleteBook))

	return router
}

func main() {
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	router := NewRouter()
	port := os.Getenv("PORT")
	log.Println(port)
	if port == "" {
		port = "3000"
	}
	log.Printf("Server running on port :%s\n", port)
	http.ListenAndServe(":"+port, router)
}

func deleteBook(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isbn := params.ByName("isbn")
	err := removeBook(isbn)
	if err != nil {
		if err == ErrBookNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func editBook(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isbn := params.ByName("isbn")
	book := new(Book)
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = updateBook(isbn, book)
	if err != nil {
		if err == ErrBookNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(DB).C(BOOKS_COLLECTION)
	index := mgo.Index{
		Key:        []string{"isbn"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}
}

func logRequest(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		log.Printf("%s %s\n", r.Method, r.URL.String())
		handler(w, r, params)
	}
}

func allBooks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	books, err := getAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(books) == 0 {
		fmt.Fprintf(w, "[]")
		return
	}
	json.NewEncoder(w).Encode(books)
}

func addBook(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	book := new(Book)
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = saveBook(book)
	if err != nil {
		if err == ErrBookAlreadyExists {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Location", r.URL.Path+"/"+book.ISBN)
}

func booksByIsbn(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	isbn := params.ByName("isbn")
	book, err := getBookByIsbn(isbn)
	if err != nil {
		if err == ErrBookNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}
