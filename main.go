package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "embed"

	"github.com/mattn/go-mastodon"
	"gopkg.in/yaml.v3"
)

//go:embed quotes.yaml
var quotesSrc []byte

func init() {
	rand.Seed(time.Now().Unix())
}

type Article struct {
	Link   string   `yaml:"link"`
	Title  string   `yaml:"title"`
	Quotes []string `yaml:"quotes"`
}

type Articles struct {
	articles []Article
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

func (a Articles) Random() Status {
	aidx := rand.Intn(len(a.articles))
	qidx := rand.Intn(len(a.articles[aidx].Quotes))

	return Status{
		Link:  a.articles[aidx].Link,
		Title: a.articles[aidx].Title,
		Quote: strings.TrimSpace(a.articles[aidx].Quotes[qidx]),
	}
}

func main() {
	ctx := context.Background()
	token := os.Getenv("MASTO_ACCESS_TOKEN")
	os.Setenv("MASTO_ACCESS_TOKEN", "<redacted>")

	var check bool
	flag.BoolVar(&check, "check", false, "check the quotes.yaml for a validity and stop")
	flag.Parse()

	as, err := UnmarshalQuotes(quotesSrc)
	if err != nil {
		log.Fatal(err)
	}

	if check {
		asLen := len(as.articles)
		quotesLen := 0
		for _, article := range as.articles {
			quotesLen += len(article.Quotes)
		}
		log.Printf("I: quotes.yaml is valid. Contains %d articles with %d quotes in total.", asLen, quotesLen)
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
