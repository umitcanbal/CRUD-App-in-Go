package main

import (
	"encoding/json"
	"fmt"
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

func createErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}

func getUsers(writer http.ResponseWriter, req *http.Request, data UsersSlice) {
	fmt.Println("getusers endpoint is requested")

	dataJson, err := json.Marshal(data)
	// err = errors.New("Something went wronggg")

	writer.Header().Set("Content-type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(createErrorResponse(err))
		return
	}

	writer.Write(dataJson)
}

func main() {
	fmt.Print("started\n")

	user1 := User{
		1,
		"nameuser1",
		30,
	}

	user2 := User{
		ID:   2,
		Name: "nameuser2",
	}

	usersSlice := UsersSlice{user1, user2}

	http.HandleFunc("/getusers", func(writer http.ResponseWriter, req *http.Request) {
		getUsers(writer, req, usersSlice)
	})
	http.ListenAndServe(":8080", nil)
	fmt.Println("server is listening")

}
