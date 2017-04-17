package main

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const usage = `
usage:
	bcrypt hash|verify <password> [<cost>] [<pass-hash>]

<password> is required for both 'hash' and 'verify'
<cost> is required only for 'hash'
<pass-hash> is required only for 'verify'
`

func main() {
	if len(os.Args) < 4 ||
		(os.Args[1] != "hash" && os.Args[1] != "verify") {
		fmt.Println(usage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	password := []byte(os.Args[2])

	switch cmd {
	case "hash":
		cost, err := strconv.Atoi(os.Args[3])
		// larger the cost factor, slower the hashing
		// so that GPU gets faster
		if err != nil {
			fmt.Println("cost msut be an integer")
			os.Exit(1)
		}
		passhash, err := bcrypt.GenerateFromPassword(password, cost)
		if err != nil {
			fmt.Printf("error hashing password: %v", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(passhash))

	case "verify":
		passhash := []byte(os.Args[3])
		err := bcrypt.CompareHashAndPassword(passhash, password)
		if err != nil {
			fmt.Println("invalid password")
		} else {
			fmt.Println("valid password :)")
		}
	}
}
