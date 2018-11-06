package main

import (
	"encoding/json"
	"net/url"
	"os"
	"text/template"
)

const refAPI = "https://api.esv.org/v3/passage/text/"

var (
	refOpts = map[string]string{
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

type esvPassage struct {
	Query       string   `json:"query"`
	Canonical   string   `json:"canonical"`
	Parsed      [][]int  `json:"parsed"`
	Passages    []string `json:"passages"`
	PassageMeta []struct {
		Canonical    string `json:"canonical"`
		PrevVerse    int    `json:"prev_verse"`
		NextVerse    int    `json:"next_verse"`
		ChapterStart []int  `json:"chapter_start"`
		ChapterEnd   []int  `json:"chapter_end"`
		PrevChapter  []int  `json:"prev_chapter"`
		NextChapter  []int  `json:"next_chapter"`
	} `json:"passage_meta"`
}

func getESVReference(query, token string) {
	vals := &url.Values{}
	vals.Set("q", query)
	// for k, v := range opts {
	// vals.Set(k, v)
	// }

	url, err := url.Parse(refAPI)
	if err != nil {
		logger.Println(err)
	}

	url.RawQuery = vals.Encode()

	body := esvRequest(url, token)

	passageResult := &esvPassage{}
	err = json.Unmarshal(body, passageResult)
	if err != nil {
		logger.Println(err)
	}

	renderReferenceResults(passageResult)
}

func renderReferenceResults(results *esvPassage) {
	t, _ := template.New("search_results").Parse("{{.}}\n")

	for _, passage := range results.Passages {
		err := t.Execute(os.Stdout, passage)
		if err != nil {
			logger.Println(err)
		}
	}
}
