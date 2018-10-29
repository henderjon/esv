package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

// passage represents the return payload for v3 of the ESV api
type passage struct {
	Detail      string   `json:"detail"`
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
	}
}

var (
	api  = "https://api.esv.org/v3/passage/text/"
	opts = map[string]string{
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

// get a passage of scripture by reference from the ESV Web API
func query(ref, token string) passage {
	var ok bool
	vals := &url.Values{}
	vals.Set("q", ref)
	// for k, v := range opts {
	// vals.Set(k, v)
	// }

	url, err := url.Parse(api)
	if err != nil {
		logger.Println(err)
	}

	url.RawQuery = vals.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)

	if len(token) == 0 {
		token, ok = os.LookupEnv("ESVTOKEN")
		if !ok {
			logger.Println("missing env var: ESVTOKEN")
		}
	}

	req.Header.Set("Authorization", "Token "+token)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		logger.Println(err)
	}

	defer res.Body.Close()

	var passage passage
	err = (json.NewDecoder(res.Body)).Decode(&passage)
	if err != nil {
		logger.Println(err)
	}

	if res.StatusCode != http.StatusOK {
		logger.Println(res.Status, passage.Detail)
	}

	return passage
}
