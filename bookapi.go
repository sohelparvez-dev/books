package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"      // PostgreSQL driver
	"github.com/google/uuid"   // UUID generator
)

// Database connection constants
const (
	host     = "localhost"
	port     = 5432
	user     = "books"
	password = "123"
	dbname   = "books"
)

// Book represents the structure of a book in the database
type Book struct {
	ID            string `json:"id"`              // Unique identifier
	Title         string `json:"title"`           // Title of the book
	Author        string `json:"author"`          // Author of the book
	PublishedYear string `json:"published_year"`  // Year of publication
	Genre         string `json:"genre"`           // Genre of the book
}

// Response structure for API messages
type Response struct {
	Message string `json:"message"` // API message response
}

var db *sql.DB // Global database connection object

// Initializes the database connection
func initDB() {
	var err error
	// Connection string for the PostgreSQL database
	dsn := "host=localhost port=5432 user=books password=123 dbname=books sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	log.Println("Database connection established")
}

// Main entry point of the application
func main() {
	initDB()         // Initialize the database connection
	defer db.Close() // Ensure the database connection is closed when the program exits

	// Register HTTP routes for different endpoints
	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/books/create", createBook)
	http.HandleFunc("/books/get/", getBookByStringID)
	http.HandleFunc("/books/update/", updateBookByID)
	http.HandleFunc("/books/delete/", deleteBookByID)

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler to fetch all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	// Query to fetch all books from the database
	rows, err := db.Query("SELECT id, title, author, publishedYear, genre FROM books")
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book // Slice to store the list of books
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear, &book.Genre); err != nil {
			http.Error(w, "Failed to parse book data", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// Respond with the list of books in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Handler to fetch a single book by ID
func getBookByStringID(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/books/get/")
	if id == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	var book Book
	// Query to fetch the book by ID
	query := "SELECT id, title, author, publishedYear, genre FROM books WHERE id = $1"
	row := db.QueryRow(query, id)
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear, &book.Genre); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the book details in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Handler to create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Validate HTTP method
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newBook Book
	// Parse the request body to extract book details
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the new book
	newBook.ID = uuid.New().String()

	// Query to insert the new book into the database
	query := `
		INSERT INTO books (id, title, author, publishedYear, genre)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.Exec(query, newBook.ID, newBook.Title, newBook.Author, newBook.PublishedYear, newBook.Genre)
	if err != nil {
		http.Error(w, "Failed to create the book", http.StatusInternalServerError)
		return
	}

	// Respond with the details of the newly created book
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBook)
}

// Handler to update an existing book by ID
func updateBookByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut { // Validate HTTP method
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the book ID from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/books/update/")
	if id == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	// Parse the request body to extract updated book details
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Query to update the book details in the database
	query := `
		UPDATE books
		SET title = $1, author = $2, publishedYear = $3, genre = $4
		WHERE id = $5
	`
	result, err := db.Exec(query, updatedBook.Title, updatedBook.Author, updatedBook.PublishedYear, updatedBook.Genre, id)
	if err != nil {
		http.Error(w, "Failed to update the book", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Book updated successfully"})
}

// Handler to delete a book by ID
func deleteBookByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete { // Validate HTTP method
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the book ID from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/books/delete/")
	if id == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	// Query to delete the book from the database
	query := "DELETE FROM books WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete the book", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Book deleted successfully"})
}

