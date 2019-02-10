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
	fmt.Printf("got user: %+v\n", user)
	fmt.Println("")

	followingMizchi, err := cli.IsFollowingUser(ctx, "mizchi")
	if err != nil {
		panic(err)
	}
	fmt.Printf("following @mizchi: %+v\n", followingMizchi)
	fmt.Println("")

	followingYaotti, err := cli.IsFollowingUser(ctx, "yaotti")
	if err != nil {
		panic(err)
	}
	fmt.Printf("following @yaotti: %+v\n", followingYaotti)
	fmt.Println("")

	usersResp, err := cli.GetUsers(ctx, 7, 50)
	if err != nil {
		panic(err)
	}
	fmt.Printf("usersResp: %+v\n", usersResp)
	fmt.Println("")
}
