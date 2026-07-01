package main

import (
	"fmt"
	"net/http"
)

type User struct {
	id   int
	name string
	age  int
}

type UsersSlice []User

func getUsers(writer http.ResponseWriter, req *http.Request, data UsersSlice) {
	fmt.Println("getusers endpoint is requested")
	fmt.Fprintln(writer, data)
}

func main() {
	fmt.Print("started\n")

	user1 := User{
		1,
		"nameuser1",
		30,
	}

	user2 := User{
		id:   2,
		name: "nameuser2",
	}

	usersSlice := UsersSlice{user1, user2}

	fmt.Printf("the value of `usersSlice` is %+v\n", usersSlice)

	http.HandleFunc("/getusers", func(writer http.ResponseWriter, req *http.Request) {
		getUsers(writer, req, usersSlice)
	})
	http.ListenAndServe(":8080", nil)
	fmt.Println("server is listening")

}
