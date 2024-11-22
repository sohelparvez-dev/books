# books
A small CRUD application using Golang backend with PostgreSQL

Systematic Design Process

Designing a meaningful REST API involves planning its architecture, endpoints, data models, and workflows. Here's a systematic design process tailored for building REST APIs:
1. Requirements Analysis
Understand the purpose of the API and the data it will handle. Let’s consider an example use case:
Scenario: A system for managing a library with the following features:
•	Managing books (CRUD operations)

2. Define the Data Models
List the core entities and their relationships. For a library, the main models could be:
1.	Book:
o	id: Unique identifier
o	title: Title of the book
o	author_id: Foreign key linking to the author
o	published_year: Year of publication
o	genre: Genre of the book

3. Design the Endpoints
2.	Structure endpoints around the resources (books) and define operations.
HTTP Method	Endpoint	Description
GET	/books	Get all books
GET	/books/{id}	Get a specific book by ID
POST	/books/create	Add a new book
PUT	/books/update/{id}	Update book details
DELETE	/books/delete/{id}	Delete a book
		
		
		

4. API Workflows
Example Workflow: Borrow a Book
1.	Client/Admin requests add new book (POST /book/createbook).
2.	If the books information is verified and saved, the client/admin gets the successfully saved data message and updated books list.
3.	Other API’s can able to Update/Delete the book information.
5. Authentication and Security
•	Use JWT (JSON Web Tokens) for authentication.
•	Protect endpoints with role-based access (e.g., only admins can create or delete books).
•	Enforce HTTPS for secure communication.
6. Response Formats
Define consistent responses with proper HTTP status codes:
•	200 OK: Successful retrieval.
•	201 Created: Successful creation.
•	400 Bad Request: Invalid input.
•	404 Not Found: Resource not found.
•	500 Internal Server Error: Unhandled errors.
Example Response: GET /books
{
  "data": [{
      "id": 1,
      "title": "Learn Go Programming",
      "author": “Mr. xyz”,
      "published_year": "2024",
      "genre": "Programming"
    }],
  "total": 1
}

7. Sample API Implementation (Books Example in Go)
package main
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)
type Book struct {
	ID            uuid    `json:"id"`
	Title         string `json:"title"`
	AuthorID      string `json:"author_id"`
	PublishedYear string `json:"published_year"`
	Genre         string `json:"genre"`
}

var books = []Book{{ID: 1, Title: "Golang in Action", AuthorID: 3, PublishedYear: "2024", Genre: "Programming"},}
func main() {
	http.HandleFunc("/books", handleBooks)
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	case http.MethodPost:
		var newBook Book
		if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		newBook.ID = len(books) + 1
		books = append(books, newBook)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBook)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

8. Database Integration
Replace the in-memory data structure with a PostgreSQL database using the steps provided earlier for connection setup.
9. API Documentation
Use tools like Swagger or Postman for documenting and testing the APIs.
10. Testing
•	Unit tests for individual handlers.
•	Integration tests to validate workflows.
•	Use tools like Postman or curl to automate endpoint testing.

This design ensures scalability, maintainability, and security.

