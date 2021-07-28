// Issuesreport prints a report of issues matching the search terms.
package main

import (
	"log"
	"os"
	"text/template"
	"time"

	"gopl.io/ch4/github"
)

// template
const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:	{{.CreatedAt | daysAgo}} days
{{end}}`

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

var report = template.Must(template.New("issuelist").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ))

func noMust() {
	report, err := template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ)
	if err != nil {
		log.Fatal(err)
	}
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

/*
Run:
$ go run gopl.io/ch4/issuesreport repo:golang/go is:open json decoder

Output:

60 issues:
----------------------------------------
Number: 33416
User:   bserdar
Title:  encoding/json: This CL adds Decoder.InternKeys
Age:    726 days
----------------------------------------
Number: 43716
User:   ggaaooppeenngg
Title:  encoding/json: increment byte counter when using decoder.Token
Age:    193 days
----------------------------------------
Number: 45628
User:   pgundlach
Title:  encoding/xml: add Decoder.InputPos
Age:    99 days
----------------------------------------
Number: 42571
User:   dsnet
Title:  encoding/json: clarify Decoder.InputOffset semantics
Age:    257 days
----------------------------------------
...
*/
