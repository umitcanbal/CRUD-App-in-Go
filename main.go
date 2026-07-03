package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID   int
	Name string
	Age  int
}

type UsersSlice []User

type ErrorResponse struct {
	Error string
}

var usersSlice UsersSlice
var db *sql.DB

func createErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}

func getUsers(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("/getusers endpoint is requested")
	writer.Header().Set("Content-type", "application/json")
	var usersSlice UsersSlice

	rows, _ := db.Query("SELECT id, name, age FROM users")
	defer rows.Close()

	for {
		if !rows.Next() {
			break
		}
		var user User
		rows.Scan(&user.ID, &user.Name, &user.Age)
		usersSlice = append(usersSlice, user)
	}

	dataJson, err := json.Marshal(usersSlice)
	// err = errors.New("Something went wronggg")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	writer.Write(dataJson)
}

func createUser(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("/createuser endpoint is requested")
	writer.Header().Set("Content-type", "application/json")

	var newUser User
	err := json.NewDecoder(req.Body).Decode(&newUser)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	if newUser.Name == "" {
		writer.WriteHeader(http.StatusBadRequest)
		err = errors.New("Failed to create the user, 'name' is not provided")
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	newUser.ID = rand.Intn(1000000)
	usersSlice = append(usersSlice, newUser)

	_ = json.NewEncoder(writer).Encode(newUser)
}

func findUserIndex(users UsersSlice, userId int) (bool, int) {
	for index, user := range users {
		if user.ID == userId {
			return true, index
		}
	}

	return false, -1
}

func deleteUser(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("/deleteuser endpoint is requested")
	writer.Header().Set("Content-Type", "application/json")

	userId := req.PathValue("userId")
	id, err := strconv.Atoi(userId)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	found, userIndex := findUserIndex(usersSlice, id)

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		err = fmt.Errorf("There is no user with id %d", id)
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	usersSlice = append(usersSlice[:userIndex], usersSlice[userIndex+1:]...)

	_ = json.NewEncoder(writer).Encode(map[string]string{"message": fmt.Sprintf("User with id %d is deleted", id)})
}

func updateUser(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var updateUser User

	err := json.NewDecoder(req.Body).Decode(&updateUser)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	found, userIndex := findUserIndex(usersSlice, updateUser.ID)
	if !found {
		writer.WriteHeader(http.StatusNotFound)
		err := fmt.Errorf("There is no user with id %d", updateUser.ID)
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	usersSlice[userIndex] = updateUser
	_ = json.NewEncoder(writer).Encode(updateUser)
}

func main() {
	fmt.Print("started\n")

	http.HandleFunc("POST /createuser", createUser)
	http.HandleFunc("DELETE /deleteuser/{userId}", deleteUser)
	http.HandleFunc("GET /getusers", getUsers)
	http.HandleFunc("PUT /updateuser", updateUser)

	var err error
	db, err = sql.Open("pgx", "postgres://localhost:5432/myapp")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = db.Ping()
	if err == nil {
		fmt.Println("DB connection is alive")
	} else {
		fmt.Println(err.Error())
		return
	}

	http.ListenAndServe(":8080", nil)
	fmt.Println("server is listening")
}
