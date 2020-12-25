package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gocolly/colly"
)

// Post is a struct representing a single post in nairaland
type Post struct {
	Author string
	URL    string
	Body   string
	Title  string
}

func main() {
	URL := "https://nairaland.com/"
	log.Println("Visiting", URL)

	c := colly.NewCollector(colly.CacheDir("./nl_cache"))
	postCollector := c.Clone()
	posts := make([]Post, 0, 100)

	c.OnHTML("td.featured.w a", func(e *colly.HTMLElement) {
		postCollector.Visit(e.Attr("href"))
	})

	postCollector.OnHTML(`table[summary="posts"] tbody`, func(e *colly.HTMLElement) {
		log.Println("Post found", e.Request.URL)
		var p Post
		p.Title = e.ChildText("tr:first-child td.bold.l.pu a[href]:nth-child(4)")
		p.Author = e.ChildText(`tr:first-child td.bold.l.pu:first-child a[class="user"]`)
		p.Body = e.ChildText("tr:nth-child(2) td.l.w.pd:first-child div.narrow")
		p.URL = e.Request.URL.String()
		posts = append(posts, p)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting: ", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited: ", r.Request.URL)
	})
	c.Visit(URL)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	//enc.SetEscapeHTML(false)
	enc.SetIndent(" ", "\t")
	err := enc.Encode(posts)
	if err != nil {
		log.Println("failed to serialize response: ", err)
		return
	}

	err = ioutil.WriteFile("nl-posts.json", buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
}
