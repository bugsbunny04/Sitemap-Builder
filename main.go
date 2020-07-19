package main

import (
	link "SITEMAP-BUILDER/linkparser"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	website := flag.String("url", "http://gophercises.com", "The url that you want to build the sitemap for")
	maxDepth := flag.Int("maxDepth", 3, "the no. of clicks deep you want to go in the website")
	flag.Parse()

	pages := bfs(*website, *maxDepth)
	toXml := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
}

func bfs(urlStr string, maxDepth int) []string {
	visited := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url, _ := range q {
			if _, ok := visited[url]; ok {
				continue
			}

			visited[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(visited))
	for url, _ := range visited {
		ret = append(ret, url)
	}
	return ret
}

func get(URLStr string) []string {
	resp, err := http.Get(URLStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)
	// fmt.Printf("%T", resp.Body)

	reqURL := resp.Request.URL
	baseURL := url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(resp.Body, base), WithPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parser(r)
	for _, x := range links {
		fmt.Println(x)
	}
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
