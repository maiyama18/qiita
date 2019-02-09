package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

// this main function works as integration test of this package
func main() {
	cli, err := New(os.Getenv("QIITA_ACCESS_TOKEN"), log.New(os.Stdout, "log", log.LstdFlags))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	user, err := cli.GetUser(ctx, "muiscript")
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Printf("got user: %+v\n", user)
	fmt.Println("")
}
