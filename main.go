package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	qiita, err := New(log.New(os.Stdout, "[LOG]", log.LstdFlags))
	if err != nil {
		panic(err)
	}

	user, err := qiita.GetUser("muiscript")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", user)
}
