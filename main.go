package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mobyvb/go-blog/blog"
)

func main() {
	input := flag.String("content", "./content", "content directory")
	output := flag.String("build", "./build", "build of html")

	flag.Parse()

	b, err := blog.BlogFromDir(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse blog: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*output, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create dir %q: %v\n", *output, err)
	}

	err = b.WriteToDir(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write blog %q: %v\n", *output, err)
	}
}
