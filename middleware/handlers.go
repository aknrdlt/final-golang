package middleware

import (
	"database/sql"
	"encoding/json"
	"final/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

func CreateBook(w http.ResponseWriter, r *http.Request) {

	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	insertID := insertBook(book)

	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	user, err := getBook(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func GetAllBook(w http.ResponseWriter, r *http.Request) {

	users, err := getAllBooks()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(users)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var user models.Book

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateBook(int64(id), user)

	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteBook(int64(id))

	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

//handler functions

func insertBook(book models.Book) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO books (title, author, year) VALUES ($1, $2, $3) RETURNING bookid`

	var id int64

	err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Year).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

func getBook(id int64) (models.Book, error) {
	db := createConnection()
	defer db.Close()

	var book models.Book

	sqlStatement := `SELECT * FROM books WHERE bookid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&book.ID, &book.Title, &book.Year, &book.Author)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return book, nil
	case nil:
		return book, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return book, err
}

func getAllBooks() ([]models.Book, error) {

	db := createConnection()

	defer db.Close()

	var books []models.Book

	sqlStatement := `SELECT * FROM books`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.Book

		err = rows.Scan(&user.ID, &user.Title, &user.Year, &user.Author)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		books = append(books, user)

	}

	return books, err
}

func updateBook(id int64, book models.Book) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `UPDATE books SET title=$2, author=$3, year=$4 WHERE bookid=$1`

	res, err := db.Exec(sqlStatement, id, book.Title, book.Author, book.Year)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func deleteBook(id int64) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM books WHERE bookid=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
