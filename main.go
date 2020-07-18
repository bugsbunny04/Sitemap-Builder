package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	website := flag.String("url", "https://gophercises.com", "The url that you want to build the sitemap for")
	flag.Parse()

	resp, err := http.Get(*website)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Body)
	io.Copy(os.Stdout, resp.Body)
}
