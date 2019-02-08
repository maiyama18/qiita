package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	qiita, err := New(log.New(os.Stdout, "[LOG]", log.LstdFlags))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	user, err := qiita.GetUser(ctx, "muiscript")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", user)
}
