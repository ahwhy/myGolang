package utils_test

import (
	"fmt"
	htmlempl "html/template"
	"log"
	"os"
	"testing"
	txtempl "text/template"
	"time"

	"github.com/ahwhy/myGolang/utils"
)

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func TestBinaryFormat(t *testing.T) {
	fmt.Println(utils.BinaryFormat(4095))
}

func TestCalculateString(t *testing.T) {
	fmt.Println(utils.CalculateString("Golang"))
}

func TestHumanBytesLoaded(t *testing.T) {
	fmt.Println(utils.HumanBytesLoaded(1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024 * 1024 * 1024))
}

func TestMd5(t *testing.T) {
	fmt.Println(utils.Md5("Atlantis"))
}

func TestMd5Salt(t *testing.T) {
	fmt.Println(utils.Md5Salt("Atlantis", ""))
}

func TestSplitMd5Salt(t *testing.T) {
	salt, md5 := utils.SplitMd5Salt(utils.Md5Salt("Atlantis", ""))
	fmt.Printf("Salt: %s \nMd5: %s\n", salt, md5)
}

func TestRandString(t *testing.T) {
	fmt.Println(utils.RandString(10))
}

func TestSearchIssuesText(t *testing.T) {
	// result, err := github.SearchIssues(os.Args[1:])
	result, err := utils.SearchIssues([]string{"repo:golang/go is:open json decoder"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	// for _, item := range result.Items {
	// 	fmt.Printf("#%-5d %9.9s %.55s\n",
	// 		item.Number, item.User.Login, item.Title)
	// }

	templ := `{{.TotalCount}} issues:
	{{range .Items}}----------------------------------------
	Number:    {{.Number}}
	User:      {{.User.Login}}
	Title:     {{.Title | printf "%.64s" }}
	Age:       {{.CreatedAt | daysAgo}} days
	{{end}}`

	var report = txtempl.Must(txtempl.New("issuelist").
		Funcs(txtempl.FuncMap{"daysAgo": daysAgo}).
		Parse(templ))

	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

func TestSearchIssuesHtml(t *testing.T) {
	// result, err := github.SearchIssues(os.Args[1:])
	result, err := utils.SearchIssues([]string{"repo:golang/go is:open json decoder"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	templ := `
	<h1>{{.TotalCount}} issues</h1>
	<table>
	<tr style='text-align: left'>
	  <th>#</th>
	  <th>State</th>
	  <th>User</th>
	  <th>Title</th>
	</tr>
	{{range .Items}}
	<tr>
	  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
	  <td>{{.State}}</td>
	  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
	  <td><a href='{{.HTMLURL}}'>{{.Title}}</td>
	</tr>
	{{end}}
	</table>
	`

	var report = htmlempl.Must(htmlempl.New("issuelist").
		Funcs(htmlempl.FuncMap{"daysAgo": daysAgo}).
		Parse(templ))

	file, err := os.Create("D:\\result.html")
	if err != nil {
		log.Fatal(err)
	}

	if err := report.Execute(file, result); err != nil {
		log.Fatal(err)
	}
}

func TestSnake(t *testing.T) {
	fmt.Println(utils.Snake("qwer1234ASDF"))
}
