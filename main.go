package main

import (
	"bufio"
	"html/template"
	"log"
	"os"
)

func main() {
	params := getOptParams()
	lookup(params.API.Reference, params.API.Token)
}

func lookup(ref, token string) {
	if len(ref) == 0 {
		stdin := bufio.NewReader(os.Stdin)
		ref, _ = stdin.ReadString('\n')
	}

	passage := query(ref, token)
	render(passage.Passages)
}

// renders the given passage reference
func render(ref []string) {
	t, _ := template.New("all").Parse("{{.}}\n")

	for _, s := range ref {
		err := t.Execute(os.Stdout, s)
		if err != nil {
			log.Fatal(err)
		}
	}
}
