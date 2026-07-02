package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
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

func createErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}

func getUsers(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("/getusers endpoint is requested")
	writer.Header().Set("Content-type", "application/json")

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

func main() {
	fmt.Print("started\n")

	http.HandleFunc("POST /createuser", createUser)
	http.HandleFunc("GET /getusers", getUsers)

	http.ListenAndServe(":8080", nil)
	fmt.Println("server is listening")
}
