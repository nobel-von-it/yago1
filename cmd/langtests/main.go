package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
}

func main() {
	var users []User
	url := "https://jsonplaceholder.typicode.com/users"

	cli := resty.New()
	_, err := cli.R().SetResult(&users).Get(url)
	if err != nil {
		panic(err)
	}

	for _, v := range users {
		fmt.Printf("User {id: %d, name: %s, email: %s}\r\n", v.ID, v.Name, v.Email)
	}
}
