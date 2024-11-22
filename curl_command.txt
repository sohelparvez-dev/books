//CURL Command from Windows:
curl http://localhost:8080/books
curl http://localhost:8080/books/get/{id}
curl -v -X POST -H "Content-Type: application/json" -d "{\"title\":\"New Book\",\"author\":\"New Author\",\"published_year\":\"2024\",\"genre\":\"ABC\"}" "http://localhost:8080/books/create"
curl -v -X PUT -H "Content-Type: application/json" -d "{\"title\":\"A Test Book\",\"author\":\"Mr. Curter\",\"published_year\":\"2022\",\"genre\":\"Test\"}" "http://localhost:8080/books/update/{id}"
curl -X DELETE http://localhost:8080/books/delete/{id}
//CURL Command from Linux:
curl http://localhost:8080/books
curl http://localhost:8080/books/get/{id}
curl -v -X POST -H "Content-Type: application/json" -d '{"title":"New Book","author":"New Author","published_year":"2024","genre":"ABC"}' http://localhost:8080/books/create
curl -v -X POST -H "Content-Type: application/json" -d '{"title":"New Book","author":"New Author","published_year":"2024","genre":"ABC"}' http://localhost:8080/books/update/{id}
curl http://localhost:8080/books/delete/{id}
