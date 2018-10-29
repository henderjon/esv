package main

import (
	"bufio"
	"os"
	"text/template"
)

func main() {
	params := getOptParams()
	pass := lookup(params.API.Reference, params.API.Token)
	render(pass.Passages)
}

func lookup(ref, token string) passage {
	if len(ref) == 0 {
		stdin := bufio.NewReader(os.Stdin)
		ref, _ = stdin.ReadString('\n')
	}

	return query(ref, token)

}

// renders the given passage reference
func render(ref []string) {
	t, _ := template.New("all").Parse("{{.}}\n")

	for _, s := range ref {
		err := t.Execute(os.Stdout, s)
		if err != nil {
			logger.Println(err)
		}
	}
}
