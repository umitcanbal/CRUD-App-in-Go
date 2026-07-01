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

type UserSlice []User

func myHttpHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("handler is executing")
	fmt.Fprintln(writer, "<h1> jaja</h1>")
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

	usersSlice := UserSlice{user1, user2}

	fmt.Printf("the value of `usersSlice` is %+v\n", usersSlice)

	http.HandleFunc("/somepath", myHttpHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("server is listening")

}
