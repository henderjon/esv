package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// passage represents the return payload for v3 of the ESV api
type passage struct {
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

	vals := &url.Values{}
	vals.Set("q", ref)
	for k, v := range opts {
		vals.Set(k, v)
	}

	url, err := url.Parse(api)
	if err != nil {
		log.Fatal(err)
	}

	url.RawQuery = vals.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)

	// token, ok := os.LookupEnv("ESVTOKEN")
	// if !ok {
	// 	log.Fatal("missing env var: ESVTOKEN")
	// }

	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Accept", "application/json")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	buf := bytes.Buffer{}

	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var passage passage
	err = json.Unmarshal(buf.Bytes(), &passage)
	if err != nil {
		log.Fatal(err)
	}

	return passage
}
