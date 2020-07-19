package main

import (
	link "SITEMAP-BUILDER/linkparser"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	website := flag.String("url", "http://gophercises.com", "The url that you want to build the sitemap for")
	flag.Parse()

	pages := get(*website)
	for _, page := range pages {
		fmt.Println(page)
	}

}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)
	// fmt.Printf("%T", resp.Body)

	reqUrl := resp.Request.URL
	baseUrl := url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(hrefs(resp.Body, base), WithPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parser(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func WithPrefix(prfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prfx)
	}
}
