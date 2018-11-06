package main

import (
	"encoding/json"
	"net/url"
	"os"
	"text/template"
)

const searchAPI = "https://api.esv.org/v3/passage/search/"

var (
	searchOpts = map[string]string{
		"include-passage-references":       "true",
		"include-first-verse-numbers":      "true",
		"include-verse-numbers":            "true",
		"include-footnotes":                "true",
		"include-footnote-body":            "true",
		"include-short-copyright":          "true",
		"include-copyright":                "false",
		"include-passage-horizontal-lines": "true",
		"include-heading-horizontal-lines": "true",
		"horizontal-line-length":           "55",
		"include-headings":                 "true",
		"include-selahs":                   "true",
		"indent-using":                     "space",
		"indent-paragraphs":                "2",
		"indent-poetry":                    "true",
		"indent-poetry-lines":              "4",
		"indent-declares":                  "40",
		"indent-psalm-doxology":            "30",
		"line-length":                      "80", // default 0 (unlimited)
	}
)

type esvSearchResults struct {
	Page         int `json:"page"`
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
	Results      []struct {
		Reference string `json:"reference"`
		Content   string `json:"content"`
	} `json:"results"`
}

func getESVSearch(query, token string) {
	vals := &url.Values{}
	vals.Set("q", query)
	// for k, v := range opts {
	// vals.Set(k, v)
	// }

	url, err := url.Parse(searchAPI)
	if err != nil {
		logger.Println(err)
	}

	url.RawQuery = vals.Encode()

	body := esvRequest(url, token)

	searchResult := &esvSearchResults{}
	err = json.Unmarshal(body, searchResult)
	if err != nil {
		logger.Println(err)
	}

	renderSearchResults(searchResult)
}

func renderSearchResults(results *esvSearchResults) {
	t := template.New("search_results")

	meta, err := t.New("meta").Parse("\n{{.TotalResults}} results; page {{.Page}} of {{.TotalPages}}\n\n")
	if err != nil {
		logger.Println(err)
	}
	references, err := t.New("results").Parse("{{.Reference}}\n{{.Content}}\n\n")
	if err != nil {
		logger.Println(err)
	}

	err = meta.Execute(os.Stdout, results)
	if err != nil {
		logger.Println(err)
	}

	for _, result := range results.Results {
		err := references.Execute(os.Stdout, result)
		if err != nil {
			logger.Println(err)
		}
	}
}
