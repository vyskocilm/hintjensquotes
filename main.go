package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	_ "embed"

	"github.com/mattn/go-mastodon"
	"gopkg.in/yaml.v3"
)

//go:embed quotes.yaml
var quotesSrc []byte

type Article struct {
	Link   string   `yaml:"link"`
	Title  string   `yaml:"title"`
	Quotes []string `yaml:"quotes"`
}

type Articles struct {
	articles []Article
}

func (a Articles) Articles() int {
	return len(a.articles)
}

func (a Articles) Quotes() int {
	acc := 0
	for _, article := range a.articles {
		acc += len(article.Quotes)
	}
	return acc
}

// StatusAt panics for wrong index
func (a Articles) StatusAt(idx int) Status {
	if idx < 0 {
		log.Fatalf("Articles.StatusAt: negative idx %d", idx)
	}
	acc := 0
	for _, article := range a.articles {
		for _, quote := range article.Quotes {
			if acc == idx {
				return Status{
					Title: article.Title,
					Link:  article.Link,
					Quote: quote,
				}
			}
			acc++
		}
	}
	log.Fatalf("Articles.StatusAt: too big idx %d requested, maximum %d", idx, a.Quotes()-1)
	panic("x")
}

func (a Articles) Random() Status {
	idx := mustRandInt(a.Quotes())

	return a.StatusAt(idx)
}

type Status struct {
	Link  string
	Title string
	Quote string
}

func (q Status) String() string {
	const template = `%s

%s
%s`
	return fmt.Sprintf(template, q.Quote, q.Title, q.Link)
}

func (q Status) Toot() string {
	const template = `%s

%s
%s
#hintjens
`
	ret := fmt.Sprintf(template, q.Quote, q.Title, q.Link)

	// the toot max size is 500 maximum
	if len(ret) > 500 {
		// len of title + link + 4 new lines + ...
		restLen := len(q.Title) + len(q.Link) + 4 + 3
		q.Quote = q.Quote[0:len(q.Quote)-restLen] + `...`
		ret = fmt.Sprintf(template, q.Quote, q.Title, q.Link)
	}
	return ret
}

func UnmarshalQuotes(b []byte) (Articles, error) {
	var quotes map[string][]Article
	err := yaml.Unmarshal(quotesSrc, &quotes)
	if err != nil {
		return Articles{}, err
	}
	return Articles{
		articles: quotes["quotes"],
	}, nil
}

func main() {
	ctx := context.Background()
	token := os.Getenv("MASTO_ACCESS_TOKEN")
	os.Setenv("MASTO_ACCESS_TOKEN", "<redacted>")

	var check bool
	flag.BoolVar(&check, "check", false, "check the quotes.yaml for a validity and stop")
	var prnt bool
	flag.BoolVar(&prnt, "print", false, "print to stdout and exit")
	flag.Parse()

	as, err := UnmarshalQuotes(quotesSrc)
	if err != nil {
		log.Fatal(err)
	}

	if check {
		log.Printf("I: quotes.yaml is valid. Contains %d articles with %d quotes in total.", as.Articles(), as.Quotes())
		return
	}

	if prnt {
		fmt.Print(as.Random().Toot())
		return
	}

	if token == "" {
		log.Fatal("MASTO_ACCESS_TOKEN is empty, please check your environment")
	}

	// post on mastodon
	c := mastodon.NewClient(&mastodon.Config{
		Server:      "https://botsin.space/",
		AccessToken: token,
	})

	status, err := c.PostStatus(ctx,
		&mastodon.Toot{
			Status:     as.Random().Toot(),
			Visibility: "public",
			Language:   "en",
		})

	if err != nil {
		log.Fatalf("c.PostStatus: %s", err)
	}
	fmt.Printf("%#v\n", status)
}

func mustRandInt(max int) int {
	x, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		log.Fatalf("crypto/rand: %s", err)
	}
	return int(x.Int64())
}
