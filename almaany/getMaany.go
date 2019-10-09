package almaany

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

type Manaa struct {
	Word         string
	WordType     string
	Explanations []string
}

const EndPoint = "https://www.almaany.com/ar/dict/ar-ar/%s/"

func ScrapePages(word string) []Manaa {
	// Remove tashkeel first
	word = removeTashkeel(word)
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf(EndPoint, word))
	if err != nil {
		return make([]Manaa, 0)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return make([]Manaa, 0)

	}

	// Load the HTML document

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return make([]Manaa, 0)
	}

	// TODO support both meaning results
	result := doc.Find(".meaning-results").First()
	// we only need the direct children
	result = result.Find(".meaning-results > li")
	// we init s with the first result
	s := result.First()
	// we first init the maany slice with the number of results.
	allMaany := make([]Manaa, result.Length())

	// loop over the length of the result
	for i := 0; i < result.Length(); i++ {
		// the name is found in span. make sure to trim it
		name := strings.TrimSpace(s.Find("span").First().Text())
		// the word type isn't inside a html element but we know that
		// it comes after the span so we can assume it's the second element
		temp := s.Contents().Get(2)
		// data would return the string representation of the element
		wordType := strings.TrimSpace(temp.Data)
		// now we need to get the list of the explanations
		explanations := s.Find("ul>li")
		// first we create a slice with the length of elements
		explanationsArray := make([]string, explanations.Length())
		explanations.Each(func(i int, r *goquery.Selection) {
			htmlRes, err := r.Html()
			if err != nil {
				log.Fatal(err)
			}
			// basic assign
			explanationsArray[i] = strings.TrimSpace(htmlRes)
		})
		wordType = strings.ReplaceAll(wordType, "(", "")
		wordType = strings.ReplaceAll(wordType, ")", "")
		wordType = strings.ReplaceAll(wordType, ":", "")
		wordType = strings.TrimSpace(wordType)

		allMaany[i] = Manaa{name, wordType, explanationsArray}
		for true {

			// sometimes we get ads instead of actual maana so we do a double check
			s = s.Next()
			if goquery.NodeName(s) != "li" && s.Nodes != nil {
				continue
			}
			break
		}

	}
	return allMaany
}
