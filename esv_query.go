package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type esvError struct {
	Detail string `json:"detail"`
}

func esvRequest(url *url.URL, token string) []byte {
	req, err := http.NewRequest("GET", url.String(), nil)

	setReqToken(req, token)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		logger.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		apiError := &esvError{}
		err := json.Unmarshal(body, apiError)
		if err != nil {
			logger.Println(err)
		}
		logger.Println(res.Status, apiError.Detail)
		return nil
	}

	return body
}

func setReqToken(req *http.Request, token string) {
	var ok bool
	if len(token) == 0 {
		token, ok = os.LookupEnv("ESVTOKEN")
		if !ok {
			logger.Println("missing env var: ESVTOKEN")
		}
	}

	req.Header.Set("Authorization", "Token "+token)
}
