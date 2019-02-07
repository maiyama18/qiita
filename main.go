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

	post, err := qiita.GetPost("a0354d9ad70c1b8225b6")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", post.ID)
	fmt.Printf("%+v\n", post.Title)
}
