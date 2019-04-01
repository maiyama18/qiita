package main

import (
	"context"
	"fmt"
	"github.com/muiscript/qiita"
	"log"
	"os"
)

func main() {
	cli, err := qiita.New(os.Getenv("QIITA_ACCESS_TOKEN"), log.New(os.Stdout, "log", log.LstdFlags))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	user, err := cli.GetUser(ctx, "yaotti")
	if err != nil {
		panic(err)
	}
	fmt.Printf("got user: %+v\n", user)
	fmt.Println("---")
	fmt.Println("ok")
	fmt.Println("---")

	myItems, err := cli.GetAuthenticatedUserItems(ctx, 2, 3)
	if err != nil {
		panic(err)
	}
	var titles []string
	for _, item := range myItems.Items {
		titles = append(titles, item.Title)
	}
	fmt.Printf("got my items: %+v\n", titles)
	fmt.Println("---")
	fmt.Println("ok")
	fmt.Println("---")
}
