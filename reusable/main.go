package main

import (
	"fmt"
	"wolung/reusable/auth"
)

func main() {
	password := "secret"
	hash, _ := auth.HashPassword(password)
	fmt.Println("Password:", password)
	fmt.Println("Hash:", hash)

	match := auth.CheckPasswordHash(password, hash)
	fmt.Println("Match: ", match)
}
