package main

import (
	"bufio"
	"os"
)

func main() {
	params := getOptParams()
	var q string

	if len(params.API.Search) != 0 {
		getESVSearch(params.API.Search, params.API.Token)
	}

	if len(params.API.Reference) != 0 {
		getESVReference(params.API.Reference, params.API.Token)
	}

	if len(params.API.Search) == 0 && len(params.API.Reference) == 0 {
		stdin := bufio.NewReader(os.Stdin)
		q, _ = stdin.ReadString('\n')
		getESVReference(q, params.API.Token)
	}
}
