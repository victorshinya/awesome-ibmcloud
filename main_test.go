package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

var (
	regexContainLink     = regexp.MustCompile(`\* \[.*\]\(.*\)`)
	regexWithDescription = regexp.MustCompile(`\* \[.*\]\(.*\) - \S.*[\.\!\?]`)
)

func TestDuplicatedContent(t *testing.T) {
	html := readMarkdown()
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		panic(err)
	}
	list := make(map[string]bool, 0)
	doc.Find("ul>li a").Each(func(i int, s *goquery.Selection) {
		t.Run(s.Text(), func(t *testing.T) {
			href, exist := s.Attr("href")
			if !exist {
				t.Error("Expected to have a href")
			}
			if list[href] {
				t.Fatalf("Duplicated link %s", href)
			}
			list[href] = true
		})
	})
}

func TestLinkFormat(t *testing.T) {
	content, err := ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}
	list := strings.Split(string(content), "\n")
	for _, item := range list {
		containLink := regexContainLink.MatchString(item)
		if containLink {
			isCorrect := regexWithDescription.MatchString(item)
			if !isCorrect {
				t.Errorf("Expected the correct usage for new links (`* [repo name](repo url) - repo description`) on %s", item)
			}
		}
	}
}

func readMarkdown() io.Reader {
	content, err := ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	md := []byte(content)
	html := markdown.ToHTML(md, parser, nil)
	r := bytes.NewReader(html)
	return r
}
