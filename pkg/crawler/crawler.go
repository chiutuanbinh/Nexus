package crawler

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type Crawler struct {
	URLfrontier []string
}

func (crawler Crawler) fetch() {
	for len(crawler.URLfrontier) > 0 {
		url := crawler.URLfrontier[0]
		res, err := http.Get(url)
		if err != nil {
			log.Printf("Fail to get %v with error %v", url, err.Error())
			continue
		}

		node, err := html.Parse(res.Body)
		if err != nil {
			log.Printf("Fail to parse body of %v with error %v", url, err)
		}

		log.Println(node)
	}
}

func Foo() {
	res, _ := http.Get("https://vnexpress.net/")

	data, _ := io.ReadAll(res.Body)
	log.Printf("%s", data)
}
