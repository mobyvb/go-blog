package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mobyvb/go-blog/blog"
)

func main() {
	title := flag.String("title", "Bloggy", "blog title")
	input := flag.String("content", "./content", "content directory")

	flag.Parse()

	b, err := blog.BlogFromDir(*title, *input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse blog: %v\n", err)
		os.Exit(1)
	}

	err = b.ServeHTTP()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start server: %v\n", err)
	}
}
