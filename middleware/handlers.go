package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"final/models"  // models package where User schema is defined
	"fmt"
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateBook create a user in the postgres db
func CreateBook(w http.ResponseWriter, r *http.Request) {

	// create an empty user of type models.User
	var book models.Book

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := insertBook(book)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetUser will return a single user by its id
func GetBook(w http.ResponseWriter, r *http.Request) {
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getUser function with user id to retrieve a single user
	user, err := getBook(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(user)
}

// GetAllUser will return all the users
func GetAllBook(w http.ResponseWriter, r *http.Request) {

	// get all the users in the db
	users, err := getAllBooks()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the users as response
	json.NewEncoder(w).Encode(users)
}

// UpdateUser update user's detail in the postgres db
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty user of type models.User
	var user models.Book

	// decode the json request to user
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update user to update the user
	updatedRows := updateBook(int64(id), user)

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteUser delete user's detail in the postgres db
func DeleteBook(w http.ResponseWriter, r *http.Request) {

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteUser, convert the int to int64
	deletedRows := deleteBook(int64(id))

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one user in the DB
func insertBook(book models.Book) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO books (title, author, year) VALUES ($1, $2, $3) RETURNING bookid`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Year).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one user from the DB by its userid
func getBook(id int64) (models.Book, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var book models.Book

	// create the select sql query
	sqlStatement := `SELECT * FROM books WHERE bookid=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
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

	// return empty user on error
	return book, err
}

// get one user from the DB by its userid
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

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE books SET title=$2, author=$3, year=$4 WHERE bookid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, book.Title, book.Author, book.Year)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete user in the DB
func deleteBook(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM books WHERE bookid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
